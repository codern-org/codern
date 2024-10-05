package middleware

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/gofiber/fiber/v2"
)

func NewScoreboardMiddleware(
	validator domain.PayloadValidator,
	authUsecase domain.AuthUsecase,
	workspaceUsecase domain.WorkspaceUsecase,
	miscUsecase domain.MiscUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		enabled, err := miscUsecase.GetFeatureFlag("scoreboard")
		if err != nil {
			return errs.New(errs.ErrGetScoreboard, "error get scoreboard ff", err)
		}
		if !enabled {
			return errs.New(errs.ErrGetScoreboard, "scoreboard feature is disabled")
		}

		var pl payload.WorkspacePath
		if ok, err := validator.Validate(&pl, ctx); !ok {
			return err
		}

		workspace, err := workspaceUsecase.GetRaw(pl.WorkspaceId)
		if err != nil {
			return errs.New(errs.SameCode, "cannot get raw workspace id %d to get scoreboard", pl.WorkspaceId, err)
		} else if workspace == nil {
			return errs.New(errs.SameCode, "workspace id %d not found", pl.WorkspaceId)
		}

		if !workspace.IsOpenScoreboard {
			sid, err := validator.ValidateAuth(ctx)
			if sid == "" {
				return err
			}
			user, err := authUsecase.Authenticate(sid)
			if err != nil {
				return errs.New(errs.SameCode, "cannot get user to get scoreboard", err)
			}

			ok, err := workspaceUsecase.HasUser(user.Id, pl.WorkspaceId)
			if !ok {
				return errs.New(errs.ErrWorkspaceNoPerm, "cannot access workspace id %d", pl.WorkspaceId)
			} else if err != nil {
				return errs.New(errs.SameCode, "cannot validate if user id %s already exist in workspace", err)
			}
		}

		return ctx.Next()
	}
}
