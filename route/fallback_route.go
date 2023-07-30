package route

import (
	"fmt"

	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
)

func ApplyFallbackRoute(app *fiber.App) {
	app.Use(func(ctx *fiber.Ctx) error {
		return ctx.
			Status(fiber.StatusNotFound).
			JSON(response.GenericResponse{
				Sucess: false,
				Errors: []response.GenericErrorResponse{
					{
						Code:    response.ErrRouteNotFound,
						Message: fmt.Sprintf("No route for %s", ctx.Path()),
					},
				},
			})
	})
}
