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

func (c *WorkspaceController) Create(ctx *fiber.Ctx) error {
	var pl payload.CreateWorkspacePayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	workspace, err := c.workspaceUsecase.Create(
		user.Id,
		&domain.CreateWorkspace{
			Name:    pl.Name,
			Profile: pl.Profile,
		},
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
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

	workspaces, err := c.workspaceUsecase.List(user.Id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspaces)
}

func (c *WorkspaceController) ListParticipant(ctx *fiber.Ctx) error {
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	participants, err := c.workspaceUsecase.ListParticipant(pl.WorkspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, participants)
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
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)
	var workspace interface{}
	var err error

	if user != nil {
		workspace, err = c.workspaceUsecase.Get(pl.WorkspaceId, user.Id)
	} else {
		workspace, err = c.workspaceUsecase.GetRaw(pl.WorkspaceId)
	}

	if err != nil {
		return err
	} else if workspace == nil {
		return errs.New(errs.ErrWorkspaceNotFound, "workspace id %d not found", pl.WorkspaceId)
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}

func (c *WorkspaceController) GetScoreboard(ctx *fiber.Ctx) error {
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	scoreboard, err := c.workspaceUsecase.GetScoreboard(pl.WorkspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, scoreboard)
}

func (c *WorkspaceController) CreateInvitation(ctx *fiber.Ctx) error {
	var pl payload.CreateInvitationPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	id, err := c.workspaceUsecase.CreateInvitation(
		pl.WorkspaceId,
		user.Id,
		pl.ValidAt,
		pl.ValidUntil,
	)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, id)
}

func (c *WorkspaceController) GetInvitations(ctx *fiber.Ctx) error {
	var pl payload.WorkspacePath
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	invitations, err := c.workspaceUsecase.GetInvitations(pl.WorkspaceId)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, invitations)
}

func (c *WorkspaceController) DeleteInvitation(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	invitationId := ctx.Params("invitationId")

	err := c.workspaceUsecase.DeleteInvitation(invitationId, user.Id)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}

func (c *WorkspaceController) JoinByInvitationCode(ctx *fiber.Ctx) error {
	user := middleware.GetUserFromCtx(ctx)
	invitationCode := ctx.Params("invitationId")

	workspace, err := c.workspaceUsecase.JoinByInvitation(user.Id, invitationCode)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, workspace)
}

func (c *WorkspaceController) Update(ctx *fiber.Ctx) error {
	var pl payload.UpdateWorkspacePayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	if err := c.workspaceUsecase.Update(
		user.Id,
		pl.WorkspaceId,
		&domain.UpdateWorkspace{
			Name:     pl.Name,
			Favorite: pl.Favorite,
			Profile:  pl.Profile,
		},
	); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, nil)
}
