package controller

import (
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/payload"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/middleware"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WorkspaceController struct {
	logger *zap.Logger

	workspaceUsecase domain.WorkspaceUsecase
}

func NewWorkspaceController(
	logger *zap.Logger,
	workspaceUsecase domain.WorkspaceUsecase,
) *WorkspaceController {
	return &WorkspaceController{
		logger:           logger,
		workspaceUsecase: workspaceUsecase,
	}
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

	workspaces, err := c.workspaceUsecase.ListFromUserId(user.Id, &domain.WorkspaceSelector{
		Participants: selector.Has("participant"),
	})
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
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
// @Param				id					path	int				true	"Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments [get]
func (c *WorkspaceController) ListAssignment(ctx *fiber.Ctx) error {
	return nil
}

// Get godoc
//
// @Summary 		Get a workspace
// @Description	Get a workspace from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				id					path	int				true	"Workspace ID"
// @Param				fields			query []string	false	"Specific fields to include in the response"	collectionFormat(csv)	Enums(participants)
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId} [get]
func (c *WorkspaceController) Get(ctx *fiber.Ctx) error {
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	selector := payload.GetFieldSelector(ctx)

	workspace, err := c.workspaceUsecase.Get(workspaceId, &domain.WorkspaceSelector{
		Participants: selector.Has("participants"),
	})
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
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
// @Param				id					path	int				true	"Workspace ID"
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspaces/{workspaceId}/assignments/{assignmentId} [get]
func (c *WorkspaceController) GetAssignment(ctx *fiber.Ctx) error {
	return nil
}
