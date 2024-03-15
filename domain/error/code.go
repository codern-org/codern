package errs

const (
	SameCode = 0

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
	ErrUpdateUser        = 2033
	ErrGoogleAuth        = 2040

	ErrGradingRequest = 4000

	ErrFilePerm = 5000

	ErrCreateUrlPath = 9000

	ErrWorkspaceNotFound          = 30000
	ErrWorkspaceNoPerm            = 30001
	ErrGetWorkspace               = 30002
	ErrListWorkspace              = 30003
	ErrWorkspaceHasUser           = 30004
	ErrWorkspaceHasAssignment     = 30005
	ErrCreateWorkspaceParticipant = 30006
	ErrUpdateWorkspaceParticipant = 30007
	ErrDeleteWorkspaceParticipant = 30008
	ErrGetRole                    = 30009
	ErrInvalidRole                = 30010
	ErrListWorkspaceParticipant   = 30011
	ErrGetScoreboard              = 30012
	ErrCreateWorkspace            = 30013
	ErrUpdateWorkspace            = 30014
	ErrDeleteWorkspace            = 30015
	ErrWorkspaceAlreadyJoin       = 30016

	ErrCreateInvitation      = 31000
	ErrGetInvitation         = 31001
	ErrDeleteInvitation      = 31002
	ErrInvitationNotFound    = 31003
	ErrInvitationNoPerm      = 31004
	ErrInvitationInvalidDate = 31005

	ErrGetAssignment        = 40000
	ErrListAssignment       = 40001
	ErrAssignmentNotFound   = 40002
	ErrAssignmentNoTestcase = 40003
	ErrCreateAssignment     = 40004
	ErrUpdateAssignment     = 40005

	ErrCreateSubmission       = 41000
	ErrCreateSubmissionResult = 41001
	ErrGetSubmission          = 41002
	ErrListSubmission         = 41003

	ErrListTestcase   = 42000
	ErrCreateTestcase = 42001
	ErrDeleteTestcase = 42002

	ErrCreateSurvey = 50000
)
