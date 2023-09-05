package usecase

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/generator"
	"github.com/codern-org/codern/platform"
	"golang.org/x/sync/errgroup"
)

type workspaceUsecase struct {
	cfgSeaweedFs        config.ConfigSeaweedFs
	seaweedfs           *platform.SeaweedFs
	rabbitMq            *platform.RabbitMq
	workspaceRepository domain.WorkspaceRepository
}

func NewWorkspaceUsecase(
	cfgSeaweedFs config.ConfigSeaweedFs,
	seaweedfs *platform.SeaweedFs,
	rabbitMq *platform.RabbitMq,
	workspaceRepository domain.WorkspaceRepository,
) domain.WorkspaceUsecase {
	return &workspaceUsecase{
		cfgSeaweedFs:        cfgSeaweedFs,
		seaweedfs:           seaweedfs,
		rabbitMq:            rabbitMq,
		workspaceRepository: workspaceRepository,
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

	var err error
	var assignment *domain.Assignment
	submission := &domain.Submission{
		Id:           id,
		AssignmentId: assignmentId,
		UserId:       userId,
		Language:     language,
		FileUrl:      filePath,
	}

	if assignment, err = u.workspaceRepository.GetAssignment(assignmentId, userId, workspaceId); err != nil {
		return domain.NewError(domain.ErrGetAssignment, "cannot get assignment")
	}

	if err = u.workspaceRepository.CreateSubmission(submission); err != nil {
		return domain.NewError(domain.ErrCreateSubmission, "cannot create submission")
	}

	// TODO: retry strategy, error
	if err = u.seaweedfs.Upload(file, 0, filePath); err != nil {
		return domain.NewError(domain.ErrFileSystem, "cannot upload file")
	}

	if submission, err = u.workspaceRepository.GetSubmission(id); err != nil {
		return domain.NewError(domain.ErrGetSubmission, "cannot get submission")
	}

	eg, egctx := errgroup.WithContext(context.Background())

	for _, testcase := range *submission.Testcases {
		testcase := testcase
		eg.Go(func() error {
			select {
			case <-egctx.Done():
				// log.Println("cancel", testcase.Id)
				return egctx.Err()
			default:
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				// log.Println("run", testcase.Id)

				// 50% chance of failure
				// if rand.Intn(2) == 1 {
				// 	log.Println("fail:", testcase.Id)
				// 	return fmt.Errorf("from run %d", testcase.Id)
				// }

				// Follow legacy version message
				return u.rabbitMq.Publish(ctx, "", "grading", false, false, map[string]interface{}{
					"id":   fmt.Sprintf("%d.%d", assignmentId, testcase.Id),
					"type": language,
					"settings": map[string]interface{}{
						"softLimitMemory": assignment.MemoryLimit,
						"softLimitTime":   assignment.TimeLimit,
					},
					"files": []string{}, // TODO: list correct files
				})
			}
		})
	}

	if err := eg.Wait(); err != nil {
		// log.Println("catch", err)
		// TODO: deal with partial outage
		return err
	}

	return nil
}

func (u *workspaceUsecase) IsUserIn(userId string, workspaceId int) (bool, error) {
	return u.workspaceRepository.IsUserIn(userId, workspaceId)
}

func (u *workspaceUsecase) IsAssignmentIn(assignmentId int, workspaceId int) (bool, error) {
	return u.workspaceRepository.IsAssignmentIn(assignmentId, workspaceId)
}

func (u *workspaceUsecase) Get(id int, selector *domain.WorkspaceSelector) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, selector)
	if workspace == nil {
		return nil, domain.NewErrorf(domain.ErrWorkspaceNotFound, "workspace id %d not found", id)
	} else if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (u *workspaceUsecase) GetAssignment(id int, userId string, workspaceId int) (*domain.Assignment, error) {
	return u.workspaceRepository.GetAssignment(id, userId, workspaceId)
}

func (u *workspaceUsecase) List(
	userId string,
	selector *domain.WorkspaceSelector,
) (*[]domain.Workspace, error) {
	return u.workspaceRepository.ListFromUserId(userId, selector)
}

func (u *workspaceUsecase) ListAssignment(userId string, workspaceId int) (*[]domain.Assignment, error) {
	return u.workspaceRepository.ListAssignment(userId, workspaceId)
}
