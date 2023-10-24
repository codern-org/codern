package controller

import (
	"os"

	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type HealthController struct{}

func NewHealthController() *HealthController {
	return &HealthController{}
}

func (c *HealthController) Check(ctx *fiber.Ctx) error {
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}
	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"hostname": hostname,
	})
}
