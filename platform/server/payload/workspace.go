package payload

type CreateSubmissionBody struct {
	Language string `form:"language" validate:"required"`
	// TODO: inspect why the file tag is not working, even if it exists.
	// SourceCode string `form:"sourcecode" validate:"required"`
}
