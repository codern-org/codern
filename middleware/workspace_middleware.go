package middleware

import (
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewWorkspaceMiddleware(
	logger *zap.Logger,
	workspaceUsecase domain.WorkspaceUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		workspaceId := ctx.Params("workspaceId")
		if workspaceId == "" {
			return ctx.Next()
		}

		user := GetUserFromCtx(ctx)
		ok, err := workspaceUsecase.IsUserIn(user.Id, workspaceId)
		if !ok {
			return response.NewErrorResponse(
				ctx,
				fiber.StatusForbidden,
				domain.NewErrorf(domain.ErrWorkspaceNoPerm, "cannot access workspace id %s", workspaceId),
			)
		} else if err != nil {
			return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}

		ctx.Locals("workspaceId", workspaceId)

		return ctx.Next()
	}
}

func GetWorkspaceIdFromCtx(ctx *fiber.Ctx) string {
	return ctx.Locals("workspaceId").(string)
}
