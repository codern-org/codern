package usecase

import (
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/generator"
)

type workspaceUsecase struct {
	workspaceRepository domain.WorkspaceRepository
	userUsecase         domain.UserUsecase
}

func NewWorkspaceUsecase(
	workspaceRepository domain.WorkspaceRepository,
	useruserUsecase domain.UserUsecase,
) domain.WorkspaceUsecase {
	return &workspaceUsecase{
		workspaceRepository: workspaceRepository,
		userUsecase:         useruserUsecase,
	}
}

func (u *workspaceUsecase) CreateInvitation(
	workspaceId int,
	inviterId string,
	validAt time.Time,
	validUntil time.Time,
) (string, error) {
	inviterRole, err := u.workspaceRepository.GetRole(inviterId, workspaceId)
	if err != nil {
		return "", errs.New(errs.SameCode, "cannot get inviter id %s role while creating invitation", inviterId, err)
	} else if inviterRole == nil {
		return "", errs.New(errs.ErrInvitationNoPerm, "cannot get role of inviter id %s", inviterId)
	} else if *inviterRole != domain.OwnerRole && *inviterRole != domain.AdminRole {
		return "", errs.New(errs.ErrInvitationNoPerm, "inviter id %s has no permission to create invitation", inviterId)
	}

	if validAt.After(validUntil) {
		return "", errs.New(errs.ErrCreateInvitation, "valid at date must be before valid until date")
	}

	var id string
	for {
		id = generator.RandStr(constant.MaxInvitationCodeChar)
		invitation, err := u.GetInvitation(id)
		if err != nil {
			return "", errs.New(errs.SameCode, "cannot get invitation to generate invitation code", err)
		} else if invitation == nil {
			break
		}
	}

	invitation := &domain.WorkspaceInvitation{
		Id:          id,
		WorkspaceId: workspaceId,
		InviterId:   inviterId,
		CreatedAt:   time.Now(),
		ValidAt:     validAt,
		ValidUntil:  validUntil,
	}

	if err = u.workspaceRepository.CreateInvitation(invitation); err != nil {
		return "", errs.New(errs.ErrCreateInvitation, "cannot create invitation", err)
	}
	return id, nil
}

func (u *workspaceUsecase) GetInvitation(id string) (*domain.WorkspaceInvitation, error) {
	invitation, err := u.workspaceRepository.GetInvitation(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetInvitation, "cannot get invitation id %s", id, err)
	}
	return invitation, nil
}

func (u *workspaceUsecase) GetInvitations(workspaceId int) ([]domain.WorkspaceInvitation, error) {
	invitations, err := u.workspaceRepository.GetInvitations(workspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrGetInvitation, "cannot get invitations in workspace id %d", workspaceId, err)
	}
	return invitations, nil
}

func (u *workspaceUsecase) DeleteInvitation(invitationId string, userId string) error {
	invitation, err := u.GetInvitation(invitationId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get invitation id %s while deleting", invitationId, err)
	} else if invitation == nil {
		return errs.New(errs.ErrInvitationNotFound, "invitation id %s not found", invitationId)
	}

	userRole, err := u.workspaceRepository.GetRole(userId, invitation.WorkspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s role while deleting", userId, err)
	} else if userRole == nil {
		return errs.New(errs.ErrInvitationNoPerm, "user id %s not found in workspace while deleting invitation", userId)
	}

	if *userRole != domain.OwnerRole && invitation.InviterId != userId {
		return errs.New(errs.ErrInvitationNoPerm, "user id %s havs no permission to delete invitation %s", userId, invitationId)
	}

	if err := u.workspaceRepository.DeleteInvitation(invitationId); err != nil {
		return errs.New(errs.ErrDeleteInvitation, "cannot delete invitation id %s", invitationId, err)
	}
	return nil
}

func (u *workspaceUsecase) CreateParticipant(workspaceId int, userId string, role domain.WorkspaceRole) error {
	user, err := u.userUsecase.Get(userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s while creating participant", userId, err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "user id %s not found while creating participant", userId)
	}

	isUserAlreadyJoined, err := u.workspaceRepository.HasUser(userId, workspaceId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot validate if user id %s already exist in workspace", userId, err)
	} else if isUserAlreadyJoined {
		return errs.New(errs.ErrWorkspaceHasUser, "user id %s is already in workspace", userId)
	}

	participant := &domain.WorkspaceParticipant{
		WorkspaceId: workspaceId,
		UserId:      userId,
		Role:        role,
		Favorite:    false,
	}

	if err := u.workspaceRepository.CreateParticipant(participant); err != nil {
		return errs.New(errs.ErrCreateWorkspaceParticipant, "cannot create participant", err)
	}
	return nil
}

func (u *workspaceUsecase) JoinByInvitation(userId string, invitationCode string) error {
	invitation, err := u.GetInvitation(invitationCode)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get invitation id %s while joining", invitationCode, err)
	} else if invitation == nil {
		return errs.New(errs.ErrInvitationNotFound, "invitation id %s not found while joining", invitationCode)
	}

	if invitation.ValidAt.After(time.Now()) {
		return errs.New(errs.ErrInvitationInvalidDate, "invitation id %s is not valid at this time yet", invitationCode)
	}
	if invitation.ValidUntil.Before(time.Now()) {
		return errs.New(errs.ErrInvitationInvalidDate, "invitation id %s is expired", invitationCode)
	}

	err = u.CreateParticipant(invitation.WorkspaceId, userId, domain.MemberRole)
	if errs.HasCode(err, errs.ErrWorkspaceHasUser) {
		return errs.New(errs.ErrWorkspaceHasUser, "user id %s is already in workspace", userId)
	} else if err != nil {
		return errs.New(errs.SameCode, "cannot create participant while joining", err)
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

func (u *workspaceUsecase) Get(id int, userId string) (*domain.Workspace, error) {
	workspace, err := u.workspaceRepository.Get(id, userId)
	if err != nil {
		return nil, errs.New(errs.ErrGetWorkspace, "cannot get workspace id %d", id, err)
	}
	if workspace != nil {
		go u.workspaceRepository.UpdateRecent(userId, workspace.Id)
	}
	return workspace, nil
}

func (u *workspaceUsecase) GetRole(
	userId string,
	workspaceId int,
) (*domain.WorkspaceRole, error) {
	userRole, err := u.workspaceRepository.GetRole(userId, workspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrGetRole, "cannot get user role", err)
	}
	return userRole, nil
}

func (u *workspaceUsecase) List(userId string) ([]domain.Workspace, error) {
	workspaces, err := u.workspaceRepository.List(userId)
	if err != nil {
		return nil, errs.New(errs.ErrListWorkspace, "cannot list workspace", err)
	}
	return workspaces, nil
}

func (u *workspaceUsecase) ListParticipant(workspaceId int) ([]domain.WorkspaceParticipant, error) {
	participants, err := u.workspaceRepository.ListParticipant(workspaceId)
	if err != nil {
		return nil, errs.New(errs.ErrListWorkspaceParticipant, "cannot list workspace particpant", err)
	}
	return participants, nil
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
