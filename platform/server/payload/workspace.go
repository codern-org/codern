package payload

import (
	"mime/multipart"

	"github.com/codern-org/codern/domain"
)

type CreateSubmissionPayload struct {
	AssignmentId int            `params:"assignmentId" validate:"required"`
	Language     string         `form:"language" validate:"required"`
	SourceCode   multipart.File `file:"sourcecode" validate:"required"`
}

type CreateWorkspacePayload struct {
	Name           string         `json:"name" validate:"required"`
	WorkspaceImage multipart.File `file:"workspace-image" validate:"required"`
}

type CreateWorkspaceParticipantPayload struct {
	UserId string `json:"user_id" validate:"required"`
}

type CreateAssignmentPayload struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description" validate:"required"`
	MemoryLimit int                    `json:"memoryLimit" validate:"required"`
	TimeLimit   int                    `json:"timeLimit" validate:"required"`
	Level       domain.AssignmentLevel `json:"level" validate:"required"`
	DetailFile  multipart.File         `file:"detail" validate:"required"`
}
