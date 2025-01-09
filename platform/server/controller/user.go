package controller

import (
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/validator"
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

func (c *UserController) Update(ctx *fiber.Ctx) error {
	var pl payload.UpdateUserPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	fileMimeType, err := validator.GetMimeType(pl.Profile)
	if err != nil {
		return err
	}
	if fileMimeType != "image/png" && fileMimeType != "image/jpeg" && fileMimeType != "image/jpg" && fileMimeType != "image/gif" {
		return errs.New(errs.ErrBodyParser, "unsupported file type")
	}

	if err := c.userUsecase.Update(
		user.Id,
		&domain.UpdateUser{
			DisplayName: pl.DisplayName,
			Profile:     pl.Profile,
		},
	); err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"updated_at": time.Now(),
	})
}

func (c *UserController) UpdatePassword(ctx *fiber.Ctx) error {
	var pl payload.UpdateUserPasswordPayload
	if ok, err := c.validator.Validate(&pl, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.userUsecase.UpdatePassword(user.Id, pl.OldPassword, pl.NewPassword)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"updated_at": time.Now(),
	})
}
