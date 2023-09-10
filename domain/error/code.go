package errs

const (
	ErrInternal   = 1
	ErrRoute      = 2
	ErrFileSystem = 3

	ErrAuthHeader      = 1000
	ErrBodyParser      = 1001
	ErrBodyValidator   = 1002
	ErrQueryParser     = 1003
	ErrQueryValidator  = 1004
	ErrParamsParser    = 1005
	ErrParamsValidator = 1006

	ErrSessionPrefix     = 2000
	ErrSignatureMismatch = 2001
	ErrInvalidSession    = 2002
	ErrSessionExpired    = 2003
	ErrInvalidEmail      = 2010
	ErrDupEmail          = 2011
	ErrUserNotFound      = 2022
	ErrUserPassword      = 2023

	ErrWorkspaceNotFound    = 3000
	ErrWorkspaceNoPerm      = 3001
	ErrCreateSubmission     = 3010
	ErrGetAssignment        = 3011
	ErrListTestcase         = 3012
	ErrAssignmentNoTestcase = 3013
)
