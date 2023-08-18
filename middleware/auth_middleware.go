package middleware

import (
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewAuthMiddleware(
	logger *zap.Logger,
	validator domain.PayloadValidator,
	authUsecase domain.AuthUsecase,
) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sid, err := validator.ValidateAuth(ctx)
		if sid == "" {
			return err
		}

		user, err := authUsecase.Authenticate(sid)
		if err != nil {
			logger.Warn(
				"Unauthorized incomming request",
				zap.String("request_id", ctx.Locals("requestid").(string)),
				zap.String("path", ctx.Path()),
				zap.String("error", err.Error()),
			)
			return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, err)
		}

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
