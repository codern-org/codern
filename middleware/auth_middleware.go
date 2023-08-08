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
			logger.Error(
				"Unauthorized incomming request",
				zap.Any("error", map[string]interface{}{
					"path":   ctx.Path(),
					"status": fiber.StatusUnauthorized,
					"error":  err.Error(),
				}),
			)
			return ctx.Status(fiber.StatusUnauthorized).JSON(response.GenericErrorResponse{
				Code:    response.ErrUnauthorized,
				Message: "Unauthorized with this auth header",
			})
		}

		ctx.Locals("user", user)

		return ctx.Next()
	}
}
