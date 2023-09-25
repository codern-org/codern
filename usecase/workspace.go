package usecase

import (
	"fmt"
	"io"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/generator"
	"github.com/codern-org/codern/platform"
)

type workspaceUsecase struct {
	seaweedfs           *platform.SeaweedFs
	workspaceRepository domain.WorkspaceRepository
	userUsecase         domain.UserUsecase
}

func NewWorkspaceUsecase(
	seaweedfs *platform.SeaweedFs,
	workspaceRepository domain.WorkspaceRepository,
	userUsecase domain.UserUsecase,
) domain.WorkspaceUsecase {
	return &workspaceUsecase{
		seaweedfs:           seaweedfs,
		workspaceRepository: workspaceRepository,
		userUsecase:         userUsecase,
	}
}

func (u *workspaceUsecase) CreateWorkspace(userId string, name string, file io.Reader) error {
	id := generator.GetId()
	filePath := fmt.Sprintf("/workspaces/%d/profile", id)

	workspace := &domain.Workspace{
		Id:         id,
		Name:       name,
		ProfileUrl: filePath,
	}

	// TODO: retry strategy, error
	if err := u.seaweedfs.Upload(file, 0, filePath); err != nil {
		return errs.New(errs.ErrFileSystem, "cannot upload file", err)
	}

	if err := u.workspaceRepository.CreateWorkspace(workspace, userId); err != nil {
		return errs.New(errs.ErrCreateWorkspace, "cannot create workspace", err)
	}

	return nil
}

func (u *workspaceUsecase) CreateParticipant(workspaceId int, userId string, role domain.WorkspaceRole) error {
	user, err := u.userUsecase.Get(userId)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot get user id %s while creating participant", userId, err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "user id %s not found while creating participant", userId)
	}

	isUserAlreadyJoined, err := u.workspaceRepository.HasUser(userId, workspaceId)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot validate if user id %s already exist in workspace", userId, err)
	} else if isUserAlreadyJoined {
		return errs.New(errs.ErrWorkspaceHasUser, "user id %s is already in workspace", userId)
	}

	participant := &domain.WorkspaceParticipant{
		WorkspaceId: workspaceId,
		UserId:      userId,
		Role:        role,
	}

	err = u.workspaceRepository.CreateParticipant(participant)
	if err != nil {
		return errs.New(errs.ErrCreateWorkspaceParticipant, "cannot create participant", err)
	}

	return nil
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
