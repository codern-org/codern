package response

import (
	errs "github.com/codern-org/codern/domain/error"
	"github.com/gofiber/fiber/v2"
)

var DomainErrCodeToHttpStatus = map[int]int{
	errs.ErrInternal:   fiber.StatusInternalServerError,
	errs.ErrRoute:      fiber.StatusNotFound,
	errs.ErrFileSystem: fiber.StatusInternalServerError,

	errs.ErrAuthHeader:       fiber.StatusBadRequest,
	errs.ErrPayloadValidator: fiber.StatusBadRequest,
	errs.ErrBodyParser:       fiber.StatusUnprocessableEntity,
	errs.ErrQueryParser:      fiber.StatusUnprocessableEntity,
	errs.ErrParamsParser:     fiber.StatusUnprocessableEntity,

	errs.ErrSessionPrefix:     fiber.StatusUnauthorized,
	errs.ErrSignatureMismatch: fiber.StatusUnauthorized,
	errs.ErrInvalidSession:    fiber.StatusUnauthorized,
	errs.ErrSessionExpired:    fiber.StatusUnauthorized,
	errs.ErrDupSession:        fiber.StatusConflict,
	errs.ErrCreateSession:     fiber.StatusInternalServerError,
	errs.ErrGetSession:        fiber.StatusInternalServerError,
	errs.ErrUnauthenticated:   fiber.StatusUnauthorized,
	errs.ErrInvalidEmail:      fiber.StatusBadRequest,
	errs.ErrDupEmail:          fiber.StatusConflict,
	errs.ErrUserPassword:      fiber.StatusUnauthorized,
	errs.ErrUserNotFound:      fiber.StatusNotFound,
	errs.ErrGetUser:           fiber.StatusInternalServerError,
	errs.ErrCreateUser:        fiber.StatusInternalServerError,
	errs.ErrGoogleAuth:        fiber.StatusInternalServerError,

	errs.ErrGradingRequest: fiber.StatusInternalServerError,

	errs.ErrFilePerm: fiber.StatusForbidden,

	errs.ErrCreateUrlPath: fiber.StatusInternalServerError,

	errs.ErrWorkspaceNotFound:          fiber.StatusNotFound,
	errs.ErrWorkspaceNoPerm:            fiber.StatusForbidden,
	errs.ErrGetWorkspace:               fiber.StatusInternalServerError,
	errs.ErrListWorkspace:              fiber.StatusInternalServerError,
	errs.ErrWorkspaceHasUser:           fiber.StatusInternalServerError,
	errs.ErrWorkspaceHasAssignment:     fiber.StatusInternalServerError,
	errs.ErrWorkspaceUpdateRole:        fiber.StatusInternalServerError,
	errs.ErrWorkspaceUpdateRolePerm:    fiber.StatusForbidden,
	errs.ErrGetRole:                    fiber.StatusInternalServerError,
	errs.ErrListWorkspaceParticipant:   fiber.StatusInternalServerError,
	errs.ErrGetScoreboard:              fiber.StatusInternalServerError,
	errs.ErrUpdateWorkspace:            fiber.StatusInternalServerError,
	errs.ErrCreateWorkspaceParticipant: fiber.StatusInternalServerError,
	errs.ErrWorkspaceAlreadyJoin:       fiber.StatusConflict,

	errs.ErrCreateInvitation:      fiber.StatusInternalServerError,
	errs.ErrGetInvitation:         fiber.StatusInternalServerError,
	errs.ErrDeleteInvitation:      fiber.StatusInternalServerError,
	errs.ErrInvitationNotFound:    fiber.StatusNotFound,
	errs.ErrInvitationNoPerm:      fiber.StatusForbidden,
	errs.ErrInvitationInvalidDate: fiber.StatusBadRequest,

	errs.ErrGetAssignment:        fiber.StatusInternalServerError,
	errs.ErrListAssignment:       fiber.StatusInternalServerError,
	errs.ErrAssignmentNotFound:   fiber.StatusNotFound,
	errs.ErrAssignmentNoTestcase: fiber.StatusInternalServerError,
	errs.ErrCreateAssignment:     fiber.StatusInternalServerError,
	errs.ErrUpdateAssignment:     fiber.StatusInternalServerError,

	errs.ErrCreateSubmission:       fiber.StatusInternalServerError,
	errs.ErrCreateSubmissionResult: fiber.StatusInternalServerError,
	errs.ErrGetSubmission:          fiber.StatusInternalServerError,
	errs.ErrListSubmission:         fiber.StatusInternalServerError,

	errs.ErrListTestcase:   fiber.StatusInternalServerError,
	errs.ErrCreateTestcase: fiber.StatusInternalServerError,
	errs.ErrDeleteTestcase: fiber.StatusInternalServerError,

	errs.ErrCreateSurvey: fiber.StatusInternalServerError,
}
