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
}

func NewAssignmentUsecase(
	seaweedfs *platform.SeaweedFs,
	assignmentRepository domain.AssignmentRepository,
	gradingPublisher domain.GradingPublisher,
) domain.AssignmentUsecase {
	return &assignmentUsecase{
		seaweedfs:            seaweedfs,
		assignmentRepository: assignmentRepository,
		gradingPublisher:     gradingPublisher,
	}
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
