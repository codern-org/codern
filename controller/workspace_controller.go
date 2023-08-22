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

// ListFromUserId godoc
//
// @Summary 		List workspaces of user
// @Description	Get all workspaces from the specific user id
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param 			userId			path	string		true "User ID"
// @Param				fields			query []string	false	"Specific fields to include in the response"	collectionFormat(csv)	Enums(ownerName,participants)
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/user/{userId}/workspace [get]
func (c *WorkspaceController) ListFromUserId(ctx *fiber.Ctx) error {
	userId, isMe := payload.GetUserIdParam(ctx)
	selector := payload.GetFieldSelector(ctx)

	if !isMe {
		return response.NewErrorResponse(
			ctx,
			fiber.StatusForbidden,
			domain.NewError(domain.ErrWorkspaceNoPerm, "Do not have permission to get a list of workspace"),
		)
	}

	workspaces, err := c.workspaceUsecase.ListFromUserId(userId, &domain.WorkspaceSelector{
		OwnerName:    selector.Has("ownerName"),
		Participants: selector.Has("participant"),
	})
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspaces)
}

// Get godoc
//
// @Summary 		Get a workspace from workspace id
// @Description	Get a workspace from workspace id on path parameter
// @Tags 				workspace
// @Accept 			json
// @Produce 		json
// @Param				id					path	int				true	"Workspace ID"
// @Param				fields			query []string	false	"Specific fields to include in the response"	collectionFormat(csv)	Enums(ownerName,participants)
// @Security 		ApiKeyAuth
// @Param 			sid header string true "Session ID"
// @Router 			/api/workspace/{id} [get]
func (c *WorkspaceController) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return response.NewErrParamResponse(ctx, "id")
	}
	selector := payload.GetFieldSelector(ctx)
	user := middleware.GetUserFromCtx(ctx)

	ok, err := c.workspaceUsecase.IsUserIn(user.Id, id)
	if !ok {
		return response.NewErrorResponse(
			ctx,
			fiber.StatusForbidden,
			domain.NewError(domain.ErrWorkspaceNoPerm, "Do not have permission to get a workspace"),
		)
	} else if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
	}

	workspace, err := c.workspaceUsecase.Get(id, &domain.WorkspaceSelector{
		OwnerName:    selector.Has("ownerName"),
		Participants: selector.Has("participants"),
	})
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}
