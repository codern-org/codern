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

	errs.ErrSessionPrefix:     401,
	errs.ErrSignatureMismatch: 401,
	errs.ErrInvalidSession:    401,
	errs.ErrSessionExpired:    401,
	errs.ErrDupSession:        409,
	errs.ErrCreateSession:     500,
	errs.ErrGetSession:        500,
	errs.ErrUnauthenticated:   401,
	errs.ErrInvalidEmail:      400,
	errs.ErrDupEmail:          409,
	errs.ErrUserPassword:      401,
	errs.ErrUserNotFound:      404,
	errs.ErrGetUser:           500,
	errs.ErrCreateUser:        500,
	errs.ErrGoogleAuth:        500,

	errs.ErrWorkspaceNotFound:       404,
	errs.ErrWorkspaceNoPerm:         403,
	errs.ErrGetWorkspace:            500,
	errs.ErrListWorkspace:           500,
	errs.ErrWorkspaceHasUser:        500,
	errs.ErrWorkspaceHasAssignment:  500,
	errs.ErrWorkspaceUpdateRole:     500,
	errs.ErrWorkspaceUpdateRolePerm: 403,
	errs.ErrCreateSubmission:        500,
	errs.ErrGetAssignment:           500,
	errs.ErrListAssignment:          500,
	errs.ErrAssignmentNotFound:      404,
	errs.ErrListTestcase:            500,
	errs.ErrAssignmentNoTestcase:    500,
	errs.ErrGetSubmission:           500,
	errs.ErrListSubmission:          500,
	errs.ErrUpdateSubmissionResult:  500,

	errs.ErrGradingRequest: 500,

	errs.ErrCreateUrlPath: 500,
}
