package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func PathType(pathType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("pathType", pathType)
		return ctx.Next()
	}
}
