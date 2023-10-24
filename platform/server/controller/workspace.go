package controller

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
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
// @Router 			/workspaces [get]
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
// @Router 			/workspaces/{workspaceId} [get]
func (c *WorkspaceController) Get(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	workspaceId := middleware.GetWorkspaceIdFromCtx(ctx)
	selector := payload.GetFieldSelector(ctx)

	workspace, err := c.workspaceUsecase.Get(
		workspaceId,
		&domain.WorkspaceSelector{Participants: selector.Has("participants")},
		user.Id,
	)
	if err != nil {
		return err
	} else if workspace == nil {
		return errs.New(errs.ErrWorkspaceNotFound, "workspace id %d not found", workspaceId)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}
