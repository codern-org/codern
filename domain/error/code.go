package errs

const (
	OverrideCode = 0

	ErrInternal   = 1
	ErrRoute      = 2
	ErrFileSystem = 3

	ErrAuthHeader       = 1000
	ErrPayloadValidator = 1001
	ErrBodyParser       = 1002
	ErrQueryParser      = 1003
	ErrParamsParser     = 1004

	ErrSessionPrefix     = 2000
	ErrSignatureMismatch = 2001
	ErrInvalidSession    = 2002
	ErrSessionExpired    = 2003
	ErrDupSession        = 2004
	ErrCreateSession     = 2005
	ErrGetSession        = 2006
	ErrUnauthenticated   = 2007
	ErrInvalidEmail      = 2010
	ErrDupEmail          = 2011
	ErrUserPassword      = 2020
	ErrUserNotFound      = 2030
	ErrGetUser           = 2031
	ErrCreateUser        = 2032
	ErrGoogleAuth        = 2040

	ErrWorkspaceNotFound       = 3000
	ErrWorkspaceNoPerm         = 3001
	ErrGetWorkspace            = 3002
	ErrListWorkspace           = 3003
	ErrWorkspaceHasUser        = 3004
	ErrWorkspaceHasAssignment  = 3005
	ErrWorkspaceUpdateRole     = 3006
	ErrWorkspaceUpdateRolePerm = 3007
	ErrCreateSubmission        = 3010
	ErrGetAssignment           = 3020
	ErrListAssignment          = 3021
	ErrAssignmentNotFound      = 3022
	ErrAssignmentNoTestcase    = 3023
	ErrListTestcase            = 3030
	ErrGetSubmission           = 3040
	ErrListSubmission          = 3041
	ErrUpdateSubmissionResult  = 3042

	ErrGradingRequest = 4000

	ErrCreateUrlPath = 9000
)
