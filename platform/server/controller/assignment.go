package controller

import (
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type AssignmentController struct {
	validator domain.PayloadValidator

	assignmentUsecase domain.AssignmentUsecase
}

func NewAssignmentController(
	validator domain.PayloadValidator,
	assignmentUsecase domain.AssignmentUsecase,
) *AssignmentController {
	return &AssignmentController{
		validator:         validator,
		assignmentUsecase: assignmentUsecase,
	}
}

// CreateSubmission godoc
//
// @Summary 		Create a new submission
// @Description	Submit a submission of the assignment
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments/{assignmentId}/submissions [post]
func (c *AssignmentController) CreateSubmission(ctx *fiber.Ctx) error {
	var pl payload.CreateSubmissionPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	err := c.assignmentUsecase.CreateSubmission(user.Id, pl.AssignmentId, workspaceId, pl.Language, pl.SourceCode)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"submitted_at": time.Now(),
	})
}

// List godoc
//
// @Summary 		List assignment
// @Description	Get all assignment from a workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments [get]
func (c *AssignmentController) List(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	assignments, err := c.assignmentUsecase.List(user.Id, workspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignments)
}

// ListSubmission godoc
//
// @Summary 		List submission
// @Description	Get all submission from a workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments/{assignmentId}/submissions [get]
func (c *AssignmentController) ListSubmission(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	submissions, err := c.assignmentUsecase.ListSubmission(user.Id, assignmentId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, submissions)
}

// Get godoc
//
// @Summary 		Get an assignment
// @Description	Get an assignment from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId					path	int				true	"Workspace ID"
// @Param				assignmentId				path	int				true	"Assignment ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments/{assignmentId} [get]
func (c *AssignmentController) Get(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	assignment, err := c.assignmentUsecase.Get(assignmentId, user.Id)
	if err != nil {
		return err
	} else if assignment == nil {
		return errs.New(errs.ErrAssignmentNotFound, "assignment id %d not found", assignmentId)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignment)
}
