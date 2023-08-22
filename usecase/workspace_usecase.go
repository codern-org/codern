package usecase

import "github.com/codern-org/codern/domain"

type workspaceUsecase struct {
	workspaceRepository domain.WorkspaceRepository
}

func NewWorkspaceUsecase(workspaceRepository domain.WorkspaceRepository) domain.WorkspaceUsecase {
	return &workspaceUsecase{workspaceRepository: workspaceRepository}
}

func (u *workspaceUsecase) CanUserView(userId string, workspaceIds []int) (bool, error) {
	return true, nil
}

func (u *workspaceUsecase) IsUserIn(userId string, workspaceId int) (bool, error) {
	return u.workspaceRepository.IsUserIn(userId, workspaceId)
}

func (u *workspaceUsecase) Get(id int, selector *domain.WorkspaceSelector) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, selector)
	if workspace == nil {
		return nil, domain.NewError(domain.ErrWorkspaceNotFound, "Requested workspace is not found")
	} else if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (u *workspaceUsecase) ListFromUserId(
	userId string,
	selector *domain.WorkspaceSelector,
) (*[]domain.Workspace, error) {
	return u.workspaceRepository.ListFromUserId(userId, selector)
}
