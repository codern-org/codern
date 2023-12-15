package payload

import (
	"time"
)

type WorkspacePath struct {
	WorkspaceId int `params:"workspaceId" validate:"required" json:"-"`
}

type AssignmentPath struct {
	WorkspacePath
	AssignmentId int `params:"assignmentId" validate:"required" json:"-"`
}

type CreateInvitationPayload struct {
	WorkspacePath
	ValidAt    time.Time `form:"validAt" validate:"required"`
	ValidUntil time.Time `form:"validUntil" validate:"required"`
}

type UpdateFavoritePayload struct {
	WorkspacePath
	Favorite bool `form:"favorite" validate:"required"`
}
