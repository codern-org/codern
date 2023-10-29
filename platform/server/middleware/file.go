package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func NewFileMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals("requestid").(string)

		if err := ctx.Next(); err != nil {
			return err
		}

		ctx.Response().Header.Set("Server", "Codern File System 1.0")
		Cors(ctx)
		ctx.Response().Header.Set(fiber.HeaderXRequestID, requestId)

		return nil
	}
}
