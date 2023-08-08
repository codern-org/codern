package response

type GenericResponse struct {
	Sucess bool                   `json:"success"`
	Errors []GenericErrorResponse `json:"errors"`
}

type GenericErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

type GenericValidationError struct {
	FailedField string
	Tag         string
	Value       interface{}
}

const (
	ErrRouteNotFound    = 1000
	ErrPayloadParser    = 1001
	ErrPayloadValidator = 1002
	ErrLoggingError     = 1003

	ErrAuthHeaderNotFound = 2000
	ErrUnauthorized       = 2001
)
