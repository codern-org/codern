package controller

import (
	"os"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"

	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		"version": constant.Version,
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

func (c *HealthController) Metrics(ctx *fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())(ctx)
}
