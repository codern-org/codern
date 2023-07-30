package response

type GenericResponse struct {
	Sucess bool                   `json:"success"`
	Errors []GenericErrorResponse `json:"errors"`
}

type GenericErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	ErrRouteNotFound = 1000
)
