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
	// TOOD: assignment validation

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

	assignment, err := u.workspaceRepository.GetAssignment(assignmentId, userId, workspaceId)
	if err != nil {
		return errs.New(errs.ErrGetAssignment, "cannot get assignment id %s", assignmentId)
	}

	testcases, err := u.workspaceRepository.ListTestcase([]int{assignmentId})
	if err != nil {
		return errs.New(errs.ErrListTestcase, "cannot list testcase")
	}
	if len(testcases) == 0 {
		return errs.New(errs.ErrAssignmentNoTestcase, "invalid assignment id %s", assignmentId)
	}

	if err := u.workspaceRepository.CreateSubmission(submission, testcases); err != nil {
		return errs.New(errs.ErrCreateSubmission, "cannot create submission")
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(file, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file")
	}

	u.gradingPublisher.Grade(assignment, submission)

	return nil
}

func (u *workspaceUsecase) IsUserIn(userId string, workspaceId int) (bool, error) {
	return u.workspaceRepository.IsUserIn(userId, workspaceId)
}

func (u *workspaceUsecase) IsAssignmentIn(assignmentId int, workspaceId int) (bool, error) {
	return u.workspaceRepository.IsAssignmentIn(assignmentId, workspaceId)
}

func (u *workspaceUsecase) Get(id int, selector *domain.WorkspaceSelector, userId string) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, selector)
	if workspace == nil {
		return nil, errs.New(errs.ErrWorkspaceNotFound, "workspace id %d not found", id)
	} else if err != nil {
		return nil, err
	}
	go u.workspaceRepository.UpdateRecent(userId, workspace.Id)
	return workspace, nil
}

func (u *workspaceUsecase) GetAssignment(id int, userId string, workspaceId int) (*domain.Assignment, error) {
	return u.workspaceRepository.GetAssignment(id, userId, workspaceId)
}

func (u *workspaceUsecase) List(
	userId string,
	selector *domain.WorkspaceSelector,
) ([]domain.Workspace, error) {
	return u.workspaceRepository.List(userId, selector)
}

func (u *workspaceUsecase) ListRecent(userId string) ([]domain.Workspace, error) {
	return u.workspaceRepository.ListRecent(userId)
}

func (u *workspaceUsecase) ListAssignment(userId string, workspaceId int) ([]domain.Assignment, error) {
	return u.workspaceRepository.ListAssignment(userId, workspaceId)
}

func (u *workspaceUsecase) ListSubmission(userId string, assignmentId int) ([]domain.Submission, error) {
	return u.workspaceRepository.ListSubmission(userId, assignmentId)
}

func (u *workspaceUsecase) UpdateSubmissionResult(result *domain.SubmissionResult) error {
	return u.workspaceRepository.UpdateSubmissionResult(result)
}
