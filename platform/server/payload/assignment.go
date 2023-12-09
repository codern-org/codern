package payload

import "mime/multipart"

type SubmissionPath struct {
	WorkspacePath
	AssignmentPath
	SubmissionId int `params:"submissionId" validate:"required" json:"-"`
}

type CreateSubmissionPayload struct {
	AssignmentPath
	Language   string         `form:"language" validate:"required"`
	SourceCode multipart.File `file:"sourcecode" validate:"required"`
}
