package payload

import "mime/multipart"

type CreateSubmissionPayload struct {
	AssignmentId int            `params:"assignmentId" validate:"required"`
	Language     string         `form:"language" validate:"required"`
	SourceCode   multipart.File `file:"sourcecode" validate:"required"`
}
