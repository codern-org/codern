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

	MaxWebSocketConnPerUser = 4
	SeaweedFsChunkSize      = 1048576 // 1 MiB
)
