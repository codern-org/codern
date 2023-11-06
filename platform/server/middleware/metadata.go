package middleware

import (
	"github.com/codern-org/codern/internal/constant"
	"github.com/gofiber/fiber/v2"
)

func PathType(pathType string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals(constant.PathTypeCtxLocal, pathType)
		return ctx.Next()
	}
}
