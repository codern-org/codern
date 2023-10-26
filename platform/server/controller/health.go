package controller

import (
	"os"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type HealthController struct {
	cfg *config.Config
}

func NewHealthController(cfg *config.Config) *HealthController {
	return &HealthController{cfg: cfg}
}

func (c *HealthController) Index(ctx *fiber.Ctx) error {
	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"name":    c.cfg.Metadata.Name,
		"version": c.cfg.Metadata.Version,
	})
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
