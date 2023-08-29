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
		if ctx.Params("workspaceId") == "" {
			return ctx.Next()
		}

		workspaceId, err := ctx.ParamsInt("workspaceId")
		if err != nil {
			return response.NewErrParamResponse(ctx, "workspaceId")
		}

		user := GetUserFromCtx(ctx)
		ok, err := workspaceUsecase.IsUserIn(user.Id, workspaceId)
		if !ok {
			return response.NewErrorResponse(
				ctx,
				fiber.StatusForbidden,
				domain.NewError(domain.ErrWorkspaceNoPerm, "Do not have permission to get a workspace"),
			)
		} else if err != nil {
			return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}

		ctx.Locals("workspaceId", workspaceId)

		return ctx.Next()
	}
}

func GetWorkspaceIdFromCtx(ctx *fiber.Ctx) int {
	return ctx.Locals("workspaceId").(int)
}
