package middleware

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/gofiber/fiber/v2"
)

func NewPublishableWorkspaceMiddleware(
	validator domain.PayloadValidator,
	authUsecase domain.AuthUsecase,
	workspaceUsecase domain.WorkspaceUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var pl payload.WorkspacePath
		if ok, err := validator.Validate(&pl, ctx); !ok {
			return err
		}

		workspace, err := workspaceUsecase.GetRaw(pl.WorkspaceId)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get raw workspace id %d", pl.WorkspaceId, err)
		} else if workspace == nil {
			return errs.New(errs.ErrWorkspaceNotFound, "workspace id %d not found", pl.WorkspaceId)
		}

		sid, err := validator.ValidateAuth(ctx)
		if !workspace.IsOpenScoreboard && sid == "" {
			return err
		}

		user, err := authUsecase.Authenticate(sid)
		if !workspace.IsOpenScoreboard && err != nil {
			return err
		}
		ctx.Locals(constant.UserCtxLocal, user)

		return ctx.Next()
	}
}

func NewWorkspaceMiddleware(
	validator domain.PayloadValidator,
	workspaceUsecase domain.WorkspaceUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		if ctx.Params("workspaceId") == "" {
			return ctx.Next()
		}

		var pl payload.WorkspacePath
		if ok, err := validator.Validate(&pl, ctx); !ok {
			return err
		}

		user := GetUserFromCtx(ctx)
		ok, err := workspaceUsecase.HasUser(user.Id, pl.WorkspaceId)
		if !ok {
			return errs.New(errs.ErrWorkspaceNoPerm, "cannot access workspace id %d", pl.WorkspaceId)
		} else if err != nil {
			return err
		}

		if ctx.Params("assignmentId") != "" {
			var pl payload.AssignmentPath
			if ok, err := validator.Validate(&pl, ctx); !ok {
				return err
			}

			ok, err := workspaceUsecase.HasAssignment(pl.AssignmentId, pl.WorkspaceId)
			if !ok {
				return errs.New(errs.ErrWorkspaceNoPerm, "cannot access assignment id %d", pl.AssignmentId)
			} else if err != nil {
				return err
			}

			if ctx.Params("testcaseFile") != "" {
				isAuthorized, err := workspaceUsecase.CheckPerm(user.Id, pl.WorkspaceId)
				if err != nil {
					return errs.New(errs.SameCode, "cannot get workspace role", err)
				}
				if !isAuthorized {
					return errs.New(errs.ErrWorkspaceNoPerm, "cannot access testcase of assignment id %d", pl.AssignmentId)
				}
			}
		}

		return ctx.Next()
	}
}
