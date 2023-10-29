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
	workspaceUsecase     domain.WorkspaceUsecase
	gradingPublisher     domain.GradingPublisher
}

func NewAssignmentUsecase(
	seaweedfs *platform.SeaweedFs,
	assignmentRepository domain.AssignmentRepository,
	workspaceUsecase domain.WorkspaceUsecase,
	gradingPublisher domain.GradingPublisher,
) domain.AssignmentUsecase {
	return &assignmentUsecase{
		seaweedfs:            seaweedfs,
		assignmentRepository: assignmentRepository,
		workspaceUsecase:     workspaceUsecase,
		gradingPublisher:     gradingPublisher,
	}
}

func (u *assignmentUsecase) CreateAssigment(
	workspaceId int,
	name string,
	description string,
	memoryLimit int,
	timeLimit int,
	level domain.AssignmentLevel,
	detailFile io.Reader,
) (*domain.Assignment, error) {
	id := generator.GetId()
	filePath := fmt.Sprintf(
		"/workspaces/%d/assignments/%d/detail/problem.md",
		workspaceId, id,
	)

	assignment := &domain.Assignment{
		Id:          id,
		WorkspaceId: workspaceId,
		Name:        name,
		Description: description,
		DetailUrl:   filePath,
		MemoryLimit: memoryLimit,
		TimeLimit:   timeLimit,
		Level:       level,
	}

	err := u.assignmentRepository.CreateAssigment(assignment)
	if err != nil {
		return nil, errs.New(errs.ErrCreateAssignment, "cannot create assignment", err)
	}

	if err := u.seaweedfs.Upload(detailFile, 0, filePath); err != nil {
		return nil, errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	return assignment, nil
}

func (u *assignmentUsecase) CreateTestcase(assignmentId int, testcaseFiles []domain.TestcaseFile) error {
	if len(testcaseFiles) == 0 {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase, testcase files is empty")
	}

	assignment, err := u.assignmentRepository.Get(assignmentId, "")
	if err != nil {
		return errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", assignmentId, err)
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

		if err := u.seaweedfs.Upload(testcaseFile.Input, 0, inputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload file", err)
		}

		if err := u.seaweedfs.Upload(testcaseFile.Output, 0, outputFilePath); err != nil {
			return errs.New(errs.ErrFileSystem, "cannot upload file", err)
		}
	}

	err = u.assignmentRepository.CreateTestcases(testcases)
	if err != nil {
		return errs.New(errs.ErrCreateTestcase, "cannot create testcase", err)
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
		return errs.New(errs.OverrideCode, "cannot get assignment id %d", assignmentId, err)
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

func (u *assignmentUsecase) Get(id int, userId string) (*domain.Assignment, error) {
	assignment, err := u.assignmentRepository.Get(id, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", id, err)
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
	assignment, err := u.Get(assignmentId, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d to list submission", assignmentId, err)
	} else if assignment == nil {
		return nil, errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", assignmentId)
	}

	role, err := u.workspaceUsecase.GetRole(userId, assignment.WorkspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrGetWorkspaceRole, "cannot get workspace role to list submission", err)
	}

	var userIdToFilter *string
	if *role != domain.OwnerRole {
		userIdToFilter = &userId
	}

	submissions, err := u.assignmentRepository.ListSubmission(&domain.SubmissionFilter{
		AssignmentId: &assignmentId,
		UserId:       userIdToFilter,
	})

	if err != nil {
		return nil, errs.New(errs.ErrListSubmission, "cannot list submission", err)
	}
	return submissions, nil
}

func (u *assignmentUsecase) UpdateSubmissionResults(
	submissionId int,
	compilationLog string,
	results []domain.SubmissionResult) error {
	err := u.assignmentRepository.UpdateSubmissionResults(submissionId, compilationLog, results)
	if err != nil {
		return errs.New(errs.ErrUpdateSubmissionResult, "cannot update submission result", err)
	}
	return nil
}
