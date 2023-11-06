package constant

import "os"

var (
	IsDevelopment = os.Getenv("ENVIRONMENT") == "development"

	SessionCookieName = "sid"

	RequestIdCtxLocal    = "requestid"
	PathTypeCtxLocal     = "pathType"
	UserCtxLocal         = "user"
	WorkspaceIdCtxLocal  = "workspaceId"
	AssignmentIdCtxLocal = "assignmentId"
)
