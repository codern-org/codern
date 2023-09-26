package payload

import "github.com/codern-org/codern/domain"

type CreateSubmissionBody struct {
	Language string `form:"language" validate:"required"`
	// TODO: inspect why the file tag is not working, even if it exists.
	// SourceCode string `form:"sourcecode" validate:"required"`
}

type CreateWorkspaceBody struct {
	Name string `json:"name" validate:"required"`
}

type CreateWorkspaceParticipantBody struct {
	UserId string               `json:"user_id" validate:"required"`
	Role   domain.WorkspaceRole `json:"role" validate:"required"`
}

type CreateAssignmentBody struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description" validate:"required"`
	MemoryLimit int                    `json:"memoryLimit" validate:"required"`
	TimeLimit   int                    `json:"timeLimit" validate:"required"`
	Level       domain.AssignmentLevel `json:"level" validate:"required"`
}
