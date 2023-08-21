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
	if err != nil {
		return nil, err
	}
	if workspace == nil {
		return nil, domain.NewGenericError(
			domain.ErrWorkspaceNotFound,
			"Requested workspace is not found",
		)
	}

	return workspace, nil
}

func (u *workspaceUsecase) GetAllFromUserId(userId string, hasParticipant bool) (*[]domain.Workspace, error) {
	return u.workspaceRepository.GetAllFromUserId(userId, hasParticipant)
}
