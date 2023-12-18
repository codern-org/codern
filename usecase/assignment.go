package usecase

import (
	"fmt"
	"io"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/generator"
	"github.com/codern-org/codern/platform"
)

type assignmentUsecase struct {
	seaweedfs            *platform.SeaweedFs
	assignmentRepository domain.AssignmentRepository
	gradingPublisher     domain.GradingPublisher
	workspaceUsecase     domain.WorkspaceUsecase
}

func NewAssignmentUsecase(
	seaweedfs *platform.SeaweedFs,
	assignmentRepository domain.AssignmentRepository,
	gradingPublisher domain.GradingPublisher,
	workspaceUsecase domain.WorkspaceUsecase,
) domain.AssignmentUsecase {
	return &assignmentUsecase{
		seaweedfs:            seaweedfs,
		assignmentRepository: assignmentRepository,
		gradingPublisher:     gradingPublisher,
		workspaceUsecase:     workspaceUsecase,
	}
}

func (u *assignmentUsecase) CreateAssignment(
	userId string,
	workspaceId int,
	name string,
	description string,
	memoryLimit int,
	timeLimit int,
	level domain.AssignmentLevel,
	detailFile io.Reader,
	testcaseFiles []domain.TestcaseFile,
) error {
	userRole, err := u.workspaceUsecase.GetRole(userId, workspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get workspace role while creating assignment", err)
	}

	if (userRole == nil) || (*userRole != domain.AdminRole && *userRole != domain.OwnerRole) {
		return errs.New(errs.ErrWorkspaceNoPerm, "permission denied")
	}

	id := generator.GetId()
	filePath := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/detail/problem.md",
		workspaceId, id,
	)

	assignment := &domain.RawAssignment{
		Id:          id,
		WorkspaceId: workspaceId,
		Name:        name,
		Description: description,
		DetailUrl:   filePath,
		MemoryLimit: memoryLimit,
		TimeLimit:   timeLimit,
		Level:       level,
	}

	if err = u.assignmentRepository.CreateAssignment(assignment); err != nil {
		return errs.New(errs.ErrCreateAssignment, "cannot create assignment", err)
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(detailFile, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	if err := u.CreateTestcase(id, testcaseFiles); err != nil {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase", err)
	}

	return nil
}

func (u *assignmentUsecase) UpdateAssignment(
	userId string,
	assignmentId int,
	name string,
	description string,
	memoryLimit int,
	timeLimit int,
	level domain.AssignmentLevel,
	detailFile io.Reader,
	testcaseFiles []domain.TestcaseFile,
) error {
	assignment, err := u.GetRaw(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get raw assignment id %d while updating assignment", assignmentId, err)
	}

	userRole, err := u.workspaceUsecase.GetRole(userId, assignment.WorkspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get workspace role while updating assignment", err)
	}

	if (userRole == nil) || (*userRole != domain.AdminRole && *userRole != domain.OwnerRole) {
		return errs.New(errs.ErrWorkspaceNoPerm, "permission denied")
	}

	detailUrl := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/detail/problem.md",
		assignment.WorkspaceId, assignmentId,
	)

	assignment.Id = assignmentId
	assignment.Name = name
	assignment.Description = description
	assignment.MemoryLimit = memoryLimit
	assignment.TimeLimit = timeLimit
	assignment.Level = level
	assignment.DetailUrl = detailUrl

	if err = u.assignmentRepository.UpdateAssignment(assignment); err != nil {
		return errs.New(errs.ErrUpdateAssignment, "cannot update assignment id %d", assignmentId, err)
	}

	// TODO: retry strategy, error
	if err = u.seaweedfs.Upload(detailFile, 0, detailUrl); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload detail file while updating assignment id %d", assignmentId, err)
	}

	if err = u.UpdateTestcases(assignmentId, testcaseFiles); err != nil {
		return errs.New(errs.ErrUpdateAssignment, "cannot update testcases by assignment id %d", assignmentId, err)
	}

	return nil
}

func (u *assignmentUsecase) CreateTestcase(assignmentId int, testcaseFiles []domain.TestcaseFile) error {
	if len(testcaseFiles) == 0 {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase, testcase files is empty")
	}

	assignment, err := u.GetRaw(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get raw assignment id %d while creating testcase", assignmentId)
	}

	testcases := make([]domain.Testcase, len(testcaseFiles))
	for i, testcaseFile := range testcaseFiles {
		id := generator.GetId()

		inputFilePath := fmt.Sprintf(
			"/workspaces/%d/assignments/%d/testcase/%d.in",
			assignment.WorkspaceId, assignmentId, i+1,
		)

		outputFilePath := fmt.Sprintf(
			"/workspaces/%d/assignments/%d/testcase/%d.out",
			assignment.WorkspaceId, assignmentId, i+1,
		)

		testcases[i] = domain.Testcase{
			Id:            id,
			AssignmentId:  assignmentId,
			InputFileUrl:  inputFilePath,
			OutputFileUrl: outputFilePath,
		}

		// TODO: retry strategy, error
		if err := u.seaweedfs.Upload(testcaseFile.Input, 0, inputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload testcase input file", err)
		}
		if err := u.seaweedfs.Upload(testcaseFile.Output, 0, outputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload testcase output file", err)
		}
	}

	if err := u.assignmentRepository.CreateTestcases(testcases); err != nil {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase", err)
	}
	return nil
}

func (u *assignmentUsecase) UpdateTestcases(assignmentId int, testcaseFiles []domain.TestcaseFile) error {
	assignment, err := u.GetRaw(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get raw assignment id %d while updating testcase", assignmentId, err)
	}

	testcaseFileUrl := fmt.Sprintf("/workspaces/%d/assignments/%d/testcase/", assignment.WorkspaceId, assignment.Id)

	if err := u.assignmentRepository.DeleteTestcasesByAssignmentId(assignmentId); err != nil {
		return errs.New(errs.ErrDeleteTestcase, "cannot delete old testcases by assignment id %d", assignmentId, err)
	}

	if err := u.seaweedfs.DeleteDirectory(testcaseFileUrl); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot delete testcase files while updating testcase by assignment id: %d", assignmentId, err)
	}

	if err := u.CreateTestcase(assignmentId, testcaseFiles); err != nil {
		return errs.New(errs.ErrCreateTestcase, "cannot create new testcase by assignment id %d", assignmentId, err)
	}

	return nil
}

func (u *assignmentUsecase) CreateSubmission(
	userId string,
	assignmentId int,
	workspaceId int,
	language string,
	file io.Reader,
) error {
	id := generator.GetId()
	filePath := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/submissions/%s/%d",
		workspaceId, assignmentId, userId, id,
	)
	submission := &domain.Submission{
		Id:           id,
		AssignmentId: assignmentId,
		UserId:       userId,
		Language:     language,
		FileUrl:      filePath,
	}

	assignment, err := u.Get(assignmentId, userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %d", assignmentId, err)
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", id)
	}

	if len(assignment.Testcases) == 0 {
		return errs.New(errs.ErrAssignmentNoTestcase, "invalid assignment id %d", assignmentId)
	}

	if err := u.assignmentRepository.CreateSubmission(submission, assignment.Testcases); err != nil {
		return errs.New(errs.ErrCreateSubmission, "cannot create submission", err)
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(file, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	// TODO: inform submission on grading publisher error
	return u.gradingPublisher.Grade(assignment, submission)
}

func (u *assignmentUsecase) CreateSubmissionResults(
	submissionId int,
	compilationLog string,
	results []domain.SubmissionResult,
) error {
	status := domain.AssignmentStatusComplete
	score := 0

	for _, result := range results {
		if result.IsPassed {
			score += 1
		} else {
			status = domain.AssignmentStatusIncompleted
		}
	}

	err := u.assignmentRepository.CreateSubmissionResults(submissionId, compilationLog, status, score, results)
	if err != nil {
		return errs.New(errs.ErrCreateSubmissionResult, "cannot update submission result", err)
	}
	return nil
}

func (u *assignmentUsecase) Get(id int, userId string) (*domain.Assignment, error) {
	assignment, err := u.assignmentRepository.Get(id, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", id, err)
	}
	return assignment, nil
}

func (u *assignmentUsecase) GetRaw(id int) (*domain.RawAssignment, error) {
	assignment, err := u.assignmentRepository.GetRaw(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get raw assignment id %d", id, err)
	}
	return assignment, nil
}

func (u *assignmentUsecase) GetSubmission(id int) (*domain.Submission, error) {
	submission, err := u.assignmentRepository.GetSubmission(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetSubmission, "cannot get submission id %d", id, err)
	}
	return submission, nil
}

func (u *assignmentUsecase) List(userId string, workspaceId int) ([]domain.Assignment, error) {
	assignments, err := u.assignmentRepository.List(userId, workspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrListAssignment, "cannot list assignment", err)
	}
	return assignments, nil
}

func (u *assignmentUsecase) ListSubmission(userId string, assignmentId int) ([]domain.Submission, error) {
	submissions, err := u.assignmentRepository.ListSubmission(userId, assignmentId)
	if err != nil {
		return nil, errs.New(errs.ErrListSubmission, "cannot list submission", err)
	}
	return submissions, nil
}
