package payload

import (
	"mime/multipart"

	"github.com/codern-org/codern/domain"
)

type SubmissionPath struct {
	AssignmentPath
	SubmissionId int `params:"submissionId" validate:"required" json:"-"`
}

type CreateSubmissionPayload struct {
	AssignmentPath
	Language   string         `form:"language" validate:"required"`
	SourceCode multipart.File `file:"sourcecode" validate:"required"`
}

type CreateAssignmentPayload struct {
	WorkspacePath
	Name                string                 `json:"name" validate:"required"`
	Description         string                 `json:"description" validate:"required"`
	MemoryLimit         int                    `json:"memoryLimit" validate:"required"`
	TimeLimit           int                    `json:"timeLimit" validate:"required"`
	Level               domain.AssignmentLevel `json:"level" validate:"required"`
	DetailFile          multipart.File         `file:"detail" validate:"required"`
	TestcaseInputFiles  []multipart.File       `file:"testcaseInput" validate:"required"`
	TestcaseOutputFiles []multipart.File       `file:"testcaseOutput" validate:"required"`
}

type UpdateAssignment struct {
	AssignmentPath
	Name                string                 `json:"name" validate:"required"`
	Description         string                 `json:"description" validate:"required"`
	MemoryLimit         int                    `json:"memoryLimit" validate:"required"`
	TimeLimit           int                    `json:"timeLimit" validate:"required"`
	Level               domain.AssignmentLevel `json:"level" validate:"required"`
	DetailFile          multipart.File         `file:"detail" validate:"required"`
	TestcaseInputFiles  []multipart.File       `file:"testcaseInput" validate:"required"`
	TestcaseOutputFiles []multipart.File       `file:"testcaseOutput" validate:"required"`
}
