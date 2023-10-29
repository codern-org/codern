package controller

import (
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	validator domain.PayloadValidator

	userUsecase domain.UserUsecase
}

func NewUserController(
	validator domain.PayloadValidator,
	userUsecase domain.UserUsecase,
) *UserController {
	return &UserController{
		validator:   validator,
		userUsecase: userUsecase,
	}
}

func (c *UserController) UpdatePassword(ctx *fiber.Ctx) error {
	var payload payload.UpdateUserPayload
	if ok, err := c.validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.userUsecase.UpdatePassword(user.Id, payload.OldPassword, payload.NewPassword)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"updated_at": time.Now(),
	})
}
