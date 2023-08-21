package controller

import (
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/middleware"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type WorkspaceController struct {
	logger    *zap.Logger
	validator domain.PayloadValidator

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

// Get godoc
//
// @Summary 		Get a workspace from workspace id
// @Description	Get a workspace from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				id					path	int			true	"Workspace ID"
// @Param				participant	query	string	false	"To show all participants information"
// @Security 		ApiKeyAuth
// @param 			sid header string true "Session ID"
// @Router 			/api/workspace/{id} [get]
func (c *WorkspaceController) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return response.NewErrParamResponse(ctx, "id")
	}
	hasParticipant := ctx.QueryBool("participant")

	workspace, err := c.workspaceUsecase.Get(id, hasParticipant)
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}

// GetAllFromUserId godoc
//
// @Summary 		Get all workspaces of user
// @Description	Get all workspaces of the authenticating user
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				participant	query	string	false	"To show all participants information"
// @Security 		ApiKeyAuth
// @param 			sid header string true "Session ID"
// @Router 			/api/workspace [get]
func (c *WorkspaceController) GetAllFromUserId(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	hasParticipant := ctx.QueryBool("participant")

	workspaces, err := c.workspaceUsecase.GetAllFromUserId(user.Id, hasParticipant)
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspaces)
}
