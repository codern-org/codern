package usecase

import "github.com/codern-org/codern/domain"

type workspaceUsecase struct {
	workspaceRepository domain.WorkspaceRepository
}

func NewWorkspaceUsecase(workspaceRepository domain.WorkspaceRepository) domain.WorkspaceUsecase {
	return &workspaceUsecase{workspaceRepository: workspaceRepository}
}

func (u *workspaceUsecase) Get(id int, hasParticipant bool) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, hasParticipant)
	if workspace == nil {
		return nil, domain.NewError(
			domain.ErrWorkspaceNotFound,
			"Requested workspace is not found",
		)
	} else if err != nil {
		return nil, err
	}
	return workspace, nil
}

func (u *workspaceUsecase) ListFromUserId(userId string, hasParticipant bool) (*[]domain.Workspace, error) {
	return u.workspaceRepository.ListFromUserId(userId, hasParticipant)
}
