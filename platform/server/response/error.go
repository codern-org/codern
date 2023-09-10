package response

import errs "github.com/codern-org/codern/domain/error"

var DomainErrCodeToHttpStatus = map[int]int{
	errs.ErrInternal:   500,
	errs.ErrRoute:      404,
	errs.ErrFileSystem: 500,

	errs.ErrAuthHeader:      400,
	errs.ErrBodyParser:      422,
	errs.ErrBodyValidator:   400,
	errs.ErrQueryParser:     422,
	errs.ErrQueryValidator:  400,
	errs.ErrParamsParser:    422,
	errs.ErrParamsValidator: 400,

	errs.ErrSessionPrefix:     400,
	errs.ErrSignatureMismatch: 400,
	errs.ErrInvalidSession:    400,
	errs.ErrSessionExpired:    401,
	errs.ErrInvalidEmail:      400,
	errs.ErrDupEmail:          409,
	errs.ErrUserNotFound:      404,
	errs.ErrUserPassword:      401,

	errs.ErrWorkspaceNotFound:    404,
	errs.ErrWorkspaceNoPerm:      403,
	errs.ErrCreateSubmission:     500,
	errs.ErrGetAssignment:        500,
	errs.ErrListTestcase:         500,
	errs.ErrAssignmentNoTestcase: 500,
}
