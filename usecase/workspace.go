package usecase

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
)

type workspaceUsecase struct {
	workspaceRepository domain.WorkspaceRepository
}

func NewWorkspaceUsecase(
	workspaceRepository domain.WorkspaceRepository,
) domain.WorkspaceUsecase {
	return &workspaceUsecase{
		workspaceRepository: workspaceRepository,
	}
}

func (u *workspaceUsecase) HasUser(userId string, workspaceId int) (bool, error) {
	isIn, err := u.workspaceRepository.HasUser(userId, workspaceId)
	if err != nil {
		return false, errs.New(errs.ErrWorkspaceHasUser, "cannot check if user is in workspace", err)
	}
	return isIn, nil
}

func (u *workspaceUsecase) HasAssignment(assignmentId int, workspaceId int) (bool, error) {
	isIn, err := u.workspaceRepository.HasAssignment(assignmentId, workspaceId)
	if err != nil {
		return false, errs.New(
			errs.ErrWorkspaceHasAssignment,
			"cannot check if assignment is in workspace",
			err,
		)
	}
	return isIn, nil
}

func (u *workspaceUsecase) Get(
	id int,
	selector *domain.WorkspaceSelector,
	userId string,
) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, selector)
	if err != nil {
		return nil, errs.New(errs.ErrGetWorkspace, "cannot get workspace id %d", id, err)
	}
	if workspace != nil {
		go u.workspaceRepository.UpdateRecent(userId, workspace.Id)
	}
	return workspace, nil
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

func (u *workspaceUsecase) UpdateRole(
	updaterUserId string,
	targetUserId string,
	workspaceId int,
	role domain.WorkspaceRole,
) error {
	updaterRole, err := u.workspaceRepository.GetRole(updaterUserId, workspaceId)
	if err != nil || updaterRole == nil {
		return errs.New(errs.ErrWorkspaceUpdateRole, "cannot get updater id %s role", updaterUserId, err)
	} else if *updaterRole == domain.OwnerRole {
		return errs.New(errs.ErrWorkspaceUpdateRolePerm, "no permission to update user role in workspace", err)
	}

	if err := u.workspaceRepository.UpdateRole(targetUserId, workspaceId, role); err != nil {
		return errs.New(errs.ErrWorkspaceUpdateRole, "cannot update role", err)
	}
	return nil
}
