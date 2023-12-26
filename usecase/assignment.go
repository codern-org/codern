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

func (u *assignmentUsecase) Create(
	userId string,
	workspaceId int,
	ca *domain.CreateAssignment,
) error {
	isAuthorized, err := u.CheckPerm(userId, workspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get workspace role while creating assignment", err)
	}
	if !isAuthorized {
		return errs.New(errs.ErrWorkspaceNoPerm, "permission denied")
	}

	id := generator.GetId()
	filePath := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/detail/problem.md",
		workspaceId, id,
	)

	assignment := &domain.Assignment{
		Id:          id,
		WorkspaceId: workspaceId,
		Name:        ca.Name,
		Description: ca.Description,
		DetailUrl:   filePath,
		MemoryLimit: ca.MemoryLimit,
		TimeLimit:   ca.TimeLimit,
		Level:       ca.Level,
	}

	if err := u.assignmentRepository.Create(assignment); err != nil {
		return errs.New(errs.ErrCreateAssignment, "cannot create assignment", err)
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(ca.DetailFile, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	if err := u.CreateTestcases(id, ca.TestcaseFiles); err != nil {
		return errs.New(errs.SameCode, "cannot create testcase while creating assignment", err)
	}

	return nil
}

func (u *assignmentUsecase) Update(
	userId string,
	assignmentId int,
	ua *domain.UpdateAssignment,
) error {
	assignment, err := u.Get(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %d while updating assignment", assignmentId, err)
	}
	if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", assignmentId)
	}

	isAuthorized, err := u.CheckPerm(userId, assignment.WorkspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get workspace role while updating assignment", err)
	}
	if !isAuthorized {
		return errs.New(errs.ErrWorkspaceNoPerm, "permission denied")
	}

	if ua.Name != nil {
		assignment.Name = *ua.Name
	}
	if ua.Description != nil {
		assignment.Description = *ua.Description
	}
	if ua.MemoryLimit != nil {
		assignment.MemoryLimit = *ua.MemoryLimit
	}
	if ua.TimeLimit != nil {
		assignment.TimeLimit = *ua.TimeLimit
	}
	if ua.Level != nil {
		assignment.Level = *ua.Level
	}

	if err := u.assignmentRepository.Update(assignment); err != nil {
		return errs.New(errs.ErrUpdateAssignment, "cannot update assignment id %d", assignmentId, err)
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(ua.DetailFile, 0, assignment.DetailUrl); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload detail file while updating assignment id %d", assignmentId, err)
	}

	if ua.TestcaseFiles != nil {
		if err := u.UpdateTestcases(assignmentId, *ua.TestcaseFiles); err != nil {
			return errs.New(errs.ErrUpdateAssignment, "cannot update testcases by assignment id %d", assignmentId, err)
		}
	}

	return nil
}

func (u *assignmentUsecase) Delete(userId string, id int) error {
	assignment, err := u.Get(id)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %d while deleting assignment", id, err)
	}
	if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", id)
	}

	isAuthorized, err := u.CheckPerm(userId, assignment.WorkspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get workspace role while deleting assignment", err)
	}
	if !isAuthorized {
		return errs.New(errs.ErrWorkspaceNoPerm, "permission denied")
	}

	if err := u.assignmentRepository.Delete(id); err != nil {
		return err
	}

	return nil
}

func (u *assignmentUsecase) CreateTestcases(assignmentId int, files []domain.TestcaseFile) error {
	if len(files) == 0 {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase, testcase files is empty")
	}

	assignment, err := u.Get(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %d while creating testcase", assignmentId)
	}

	testcases := make([]domain.Testcase, len(files))
	for i, file := range files {
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
		if err := u.seaweedfs.Upload(file.Input, 0, inputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload testcase input file", err)
		}
		if err := u.seaweedfs.Upload(file.Output, 0, outputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload testcase output file", err)
		}
	}

	if err := u.assignmentRepository.CreateTestcases(testcases); err != nil {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase", err)
	}
	return nil
}

func (u *assignmentUsecase) UpdateTestcases(assignmentId int, testcaseFiles []domain.TestcaseFile) error {
	assignment, err := u.Get(assignmentId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get assignment id %d while updating testcase", assignmentId, err)
	}

	testcaseFileUrl := fmt.Sprintf("/workspaces/%d/assignments/%d/testcase/", assignment.WorkspaceId, assignment.Id)

	if err := u.assignmentRepository.DeleteTestcases(assignmentId); err != nil {
		return errs.New(errs.ErrDeleteTestcase, "cannot delete old testcases by assignment id %d", assignmentId, err)
	}

	if err := u.seaweedfs.DeleteDirectory(testcaseFileUrl); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot delete testcase files while updating testcase by assignment id: %d", assignmentId, err)
	}

	if err := u.CreateTestcases(assignmentId, testcaseFiles); err != nil {
		return errs.New(errs.SameCode, "cannot create new testcase by assignment id %d", assignmentId, err)
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

	assignment, err := u.GetWithStatus(assignmentId, userId)
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

func (u *assignmentUsecase) Get(id int) (*domain.Assignment, error) {
	assignment, err := u.assignmentRepository.Get(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", id, err)
	}
	return assignment, nil
}

func (u *assignmentUsecase) GetWithStatus(id int, userId string) (*domain.AssignmentWithStatus, error) {
	assignment, err := u.assignmentRepository.GetWithStatus(id, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", id, err)
	}
	return assignment, nil
}

func (u *assignmentUsecase) CheckPerm(userId string, workspaceId int) (bool, error) {
	userRole, err := u.workspaceUsecase.GetRole(userId, workspaceId)
	if err != nil {
		return false, errs.New(errs.SameCode, "cannot get workspace role", err)
	}
	return ((userRole != nil) && (*userRole == domain.AdminRole || *userRole == domain.OwnerRole)), nil
}

func (u *assignmentUsecase) GetSubmission(id int) (*domain.Submission, error) {
	submission, err := u.assignmentRepository.GetSubmission(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetSubmission, "cannot get submission id %d", id, err)
	}
	return submission, nil
}

func (u *assignmentUsecase) List(userId string, workspaceId int) ([]domain.AssignmentWithStatus, error) {
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
