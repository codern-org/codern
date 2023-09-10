package middleware

import (
	"github.com/codern-org/codern/domain"
	"github.com/gofiber/fiber/v2"
)

func NewAuthMiddleware(
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
			return err
		}

		ctx.Locals("user", user)

		return ctx.Next()
	}
}

func GetUserFromCtx(ctx *fiber.Ctx) *domain.User {
	return ctx.Locals("user").(*domain.User)
}
