package payload

import (
	"mime/multipart"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
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
	PublishDate         time.Time              `json:"publishDate" validate:"required"`
	DueDate             *time.Time             `json:"dueDate"`
	DetailFile          multipart.File         `file:"detail" validate:"required"`
	TestcaseInputFiles  []multipart.File       `file:"testcaseInput" validate:"required"`
	TestcaseOutputFiles []multipart.File       `file:"testcaseOutput" validate:"required"`
}

type UpdateAssignment struct {
	AssignmentPath
	Name                *string                 `json:"name"`
	Description         *string                 `json:"description"`
	MemoryLimit         *int                    `json:"memoryLimit"`
	TimeLimit           *int                    `json:"timeLimit"`
	Level               *domain.AssignmentLevel `json:"level"`
	PublishDate         *time.Time              `json:"publishDate"`
	DueDate             *time.Time              `json:"dueDate"`
	DetailFile          multipart.File          `file:"detail"`
	TestcaseInputFiles  []multipart.File        `file:"testcaseInput"`
	TestcaseOutputFiles []multipart.File        `file:"testcaseOutput"`
}

type DeleteAssignment struct {
	AssignmentPath
}

func ValidateTestcaseFiles(inputs []multipart.File, outputs []multipart.File) error {
	if len(inputs) != len(outputs) {
		return errs.NewPayloadError([]errs.ValidationErrorDetail{
			{
				Field: "TestcaseInputFiles",
				Type:  "length_mismatch",
			},
			{
				Field: "TestcaseOutputFiles",
				Type:  "length_mismatch",
			},
		})
	}
	return nil
}
