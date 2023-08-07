package route

import (
	_ "github.com/codern-org/codern/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func ApplySwaggerRoutes(app *fiber.App) {
	app.Get("/swagger/*", swagger.HandlerDefault)
}
