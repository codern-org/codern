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
			return response.NewErrorResponse(
				ctx,
				fiber.StatusBadRequest,
				domain.NewErrorWithData(domain.ErrParamsValidator, "params payload is invalid", err),
			)
		}

		user := GetUserFromCtx(ctx)
		ok, err := workspaceUsecase.IsUserIn(user.Id, workspaceId)
		if !ok {
			return response.NewErrorResponse(
				ctx,
				fiber.StatusForbidden,
				domain.NewErrorf(domain.ErrWorkspaceNoPerm, "cannot access workspace id %d", workspaceId),
			)
		} else if err != nil {
			return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
		}

		var assignmentId int
		if ctx.Params("assignmentId") != "" {
			assignmentId, err = ctx.ParamsInt("assignmentId")
			if err != nil {
				return response.NewErrorResponse(
					ctx,
					fiber.StatusBadRequest,
					domain.NewErrorWithData(domain.ErrParamsValidator, "params payload is invalid", err),
				)
			}

			ok, err := workspaceUsecase.IsAssignmentIn(assignmentId, workspaceId)

			if !ok {
				return response.NewErrorResponse(
					ctx,
					fiber.StatusForbidden,
					domain.NewErrorf(domain.ErrWorkspaceNoPerm, "cannot access assignment id %d", assignmentId),
				)
			} else if err != nil {
				return response.NewErrorResponse(ctx, fiber.StatusInternalServerError, err)
			}
		}

		ctx.Locals("workspaceId", workspaceId)
		ctx.Locals("assignmentId", assignmentId)

		return ctx.Next()
	}
}

func GetWorkspaceIdFromCtx(ctx *fiber.Ctx) int {
	return ctx.Locals("workspaceId").(int)
}

func GetAssignmentIdFromCtx(ctx *fiber.Ctx) int {
	return ctx.Locals("assignmentId").(int)
}
