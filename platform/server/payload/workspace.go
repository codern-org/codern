package payload

import (
	"mime/multipart"
	"time"
)

type CreateSubmissionPayload struct {
	AssignmentId int            `params:"assignmentId" validate:"required"`
	Language     string         `form:"language" validate:"required"`
	SourceCode   multipart.File `file:"sourcecode" validate:"required"`
}

type CreateInvitationPayload struct {
	ValidAt    time.Time `form:"validAt" validate:"required"`
	ValidUntil time.Time `form:"validUntil" validate:"required"`
}
