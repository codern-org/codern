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
	ValidAt    time.Time `json:"validAt" validate:"required"`
	ValidUntil time.Time `json:"validUntil" validate:"required"`
}

type UpdateWorkspacePayload struct {
	WorkspacePath
	Name     *string `json:"name"`
	Favorite *bool   `json:"favorite"`
}
