package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func NewFileMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		requestId := ctx.Locals("requestid").(string)
		chainErr := ctx.Next()

		// Manually call error handler
		if chainErr != nil {
			if err := ctx.App().ErrorHandler(ctx, chainErr); err != nil {
				ctx.SendStatus(fiber.StatusInternalServerError)
			}
		}

		ctx.Response().Header.Set("Server", "Codern File System 1.0")
		Cors(ctx)
		ctx.Response().Header.Set(fiber.HeaderXRequestID, requestId)

		return nil
	}
}
