package errs

type ValidationErrorDetail struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

func NewValidationErr(code int, message string, details []ValidationErrorDetail) *DomainError {
	return &DomainError{
		Code:    code,
		Message: message,
		Details: details,
	}
}
