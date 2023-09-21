package usecase

import (
	"fmt"
	"io"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/generator"
	"github.com/codern-org/codern/platform"
)

type workspaceUsecase struct {
	cfg                 *config.Config
	seaweedfs           *platform.SeaweedFs
	rabbitMq            *platform.RabbitMq
	workspaceRepository domain.WorkspaceRepository
	gradingPublisher    domain.GradingPublisher
}

func NewWorkspaceUsecase(
	cfg *config.Config,
	seaweedfs *platform.SeaweedFs,
	rabbitMq *platform.RabbitMq,
	workspaceRepository domain.WorkspaceRepository,
	gradingPublisher domain.GradingPublisher,
) domain.WorkspaceUsecase {
	return &workspaceUsecase{
		cfg:                 cfg,
		seaweedfs:           seaweedfs,
		rabbitMq:            rabbitMq,
		workspaceRepository: workspaceRepository,
		gradingPublisher:    gradingPublisher,
	}
}

func (u *workspaceUsecase) CreateSubmission(
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

	assignment, err := u.GetAssignment(assignmentId, userId)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot get assignment id %d", assignmentId, err)
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", id)
	}

	if len(assignment.Testcases) == 0 {
		return errs.New(errs.ErrAssignmentNoTestcase, "invalid assignment id %d", assignmentId)
	}

	if err := u.workspaceRepository.CreateSubmission(submission, assignment.Testcases); err != nil {
		return errs.New(errs.ErrCreateSubmission, "cannot create submission", err)
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(file, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	return u.gradingPublisher.Grade(assignment, submission)
}

func (u *workspaceUsecase) IsUserIn(userId string, workspaceId int) (bool, error) {
	isIn, err := u.workspaceRepository.IsUserIn(userId, workspaceId)
	if err != nil {
		return false, errs.New(errs.ErrIsUserIn, "cannot check if user is in workspace", err)
	}
	return isIn, nil
}

func (u *workspaceUsecase) IsAssignmentIn(assignmentId int, workspaceId int) (bool, error) {
	isIn, err := u.workspaceRepository.IsAssignmentIn(assignmentId, workspaceId)
	if err != nil {
		return false, errs.New(errs.ErrIsAssignmentIn, "cannot check if assignment is in workspace", err)
	}
	return isIn, nil
}

func (u *workspaceUsecase) Get(id int, selector *domain.WorkspaceSelector, userId string) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, selector)
	if err != nil {
		return nil, errs.New(errs.ErrGetWorkspace, "cannot get workspace id %d", id, err)
	}
	if workspace != nil {
		go u.workspaceRepository.UpdateRecent(userId, workspace.Id)
	}
	return workspace, nil
}

func (u *workspaceUsecase) GetAssignment(id int, userId string) (*domain.Assignment, error) {
	assignment, err := u.workspaceRepository.GetAssignment(id, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetAssignment, "cannot get assignment id %d", id, err)
	}
	return assignment, nil
}

func (u *workspaceUsecase) GetSubmission(id int) (*domain.Submission, error) {
	submission, err := u.workspaceRepository.GetSubmission(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetSubmission, "cannot get submission id %d", id, err)
	}
	return submission, nil
}

func (u *workspaceUsecase) List(
	userId string,
	selector *domain.WorkspaceSelector,
) ([]domain.Workspace, error) {
	workspaces, err := u.workspaceRepository.List(userId, selector)
	if err != nil {
		return nil, errs.New(errs.ErrListWorkspace, "cannot list workspace", err)
	}
	return workspaces, nil
}

func (u *workspaceUsecase) ListRecent(userId string) ([]domain.Workspace, error) {
	workspaces, err := u.workspaceRepository.ListRecent(userId)
	if err != nil {
		return nil, errs.New(errs.ErrListWorkspace, "cannot list recent workspace", err)
	}
	return workspaces, nil
}

func (u *workspaceUsecase) ListAssignment(userId string, workspaceId int) ([]domain.Assignment, error) {
	assignments, err := u.workspaceRepository.ListAssignment(userId, workspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrListAssignment, "cannot list assignment", err)
	}
	return assignments, nil
}

func (u *workspaceUsecase) ListSubmission(userId string, assignmentId int) ([]domain.Submission, error) {
	submissions, err := u.workspaceRepository.ListSubmission(userId, assignmentId)
	if err != nil {
		return nil, errs.New(errs.ErrListSubmission, "cannot list submission", err)
	}
	return submissions, nil
}

func (u *workspaceUsecase) UpdateSubmissionResults(
	submissionId int,
	compilationLog string,
	results []domain.SubmissionResult) error {
	err := u.workspaceRepository.UpdateSubmissionResults(submissionId, compilationLog, results)
	if err != nil {
		return errs.New(errs.ErrUpdateSubmissionResult, "cannot update submission result", err)
	}
	return nil
}
