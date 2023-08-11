package route

import (
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
)

func ApplyFallbackRoute(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		return response.NewErrorResponse(
			ctx,
			fiber.StatusNotFound,
			domain.NewGenericError(domain.ErrRoute, fmt.Sprintf("No route for %s", ctx.Path())))
	})
}
