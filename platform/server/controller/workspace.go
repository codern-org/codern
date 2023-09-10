package controller

import (
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type WorkspaceController struct {
	validator domain.PayloadValidator

	workspaceUsecase domain.WorkspaceUsecase
}

func NewWorkspaceController(
	validator domain.PayloadValidator,
	workspaceUsecase domain.WorkspaceUsecase,
) *WorkspaceController {
	return &WorkspaceController{
		validator:        validator,
		workspaceUsecase: workspaceUsecase,
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
func (c *WorkspaceController) CreateSubmission(ctx *fiber.Ctx) error {
	var body payload.CreateSubmissionBody
	if ok, err := c.validator.ValidateBody(&body, ctx); !ok {
		return err
	}
	file, err := payload.GetFile("sourcecode", ctx)
	if err != nil {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	assignmentId, _ := ctx.ParamsInt("assignmentId")

	err = c.workspaceUsecase.CreateSubmission(user.Id, assignmentId, workspaceId, body.Language, file)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"submitted_at": time.Now(),
	})
}

// List godoc
//
// @Summary 		List workspaces
// @Description	Get all workspaces
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				fields			query []string	false	"Specific fields to include in the response"	collectionFormat(csv)	Enums(participants)
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces [get]
func (c *WorkspaceController) List(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	selector := payload.GetFieldSelector(ctx)
	order := ctx.Query("order")

	var workspaces []domain.Workspace
	var err error

	if order == "recent" {
		workspaces, err = c.workspaceUsecase.ListRecent(user.Id)
	} else {
		workspaces, err = c.workspaceUsecase.List(user.Id, &domain.WorkspaceSelector{
			Participants: selector.Has("participants"),
		})
	}
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspaces)
}

// ListAssignment godoc
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
func (c *WorkspaceController) ListAssignment(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)

	assignments, err := c.workspaceUsecase.ListAssignment(user.Id, workspaceId)
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
func (c *WorkspaceController) ListSubmission(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	submissions, err := c.workspaceUsecase.ListSubmission(user.Id, assignmentId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, submissions)
}

// Get godoc
//
// @Summary 		Get a workspace
// @Description	Get a workspace from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				workspaceId	path	int				true	"Workspace ID"
// @Param				fields			query []string	false	"Specific fields to include in the response"	collectionFormat(csv)	Enums(participants)
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId} [get]
func (c *WorkspaceController) Get(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	selector := payload.GetFieldSelector(ctx)

	workspace, err := c.workspaceUsecase.Get(workspaceId, &domain.WorkspaceSelector{
		Participants: selector.Has("participants"),
	}, user.Id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}

// GetAssignment godoc
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
func (c *WorkspaceController) GetAssignment(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	assignmentId := middleware.GetAssignmentIdFromCtx(ctx)

	assignment, err := c.workspaceUsecase.GetAssignment(assignmentId, user.Id, workspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, assignment)
}
