package response

import (
	"errors"

	"github.com/codern-org/codern/domain"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	var domainError *domain.Error
	if errors.As(err, &domainError) {
		return ctx.Status(status).JSON(Response{
			Success: false,
			Error:   err,
		})
	}

	ctx.Locals("internal_error", err)
	return ctx.Status(fiber.StatusInternalServerError).JSON(Response{
		Success: false,
		Error: &domain.Error{
			Code:    domain.ErrInternal,
			Message: err.Error(),
		},
	})
}

func NewSuccessResponse(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(Response{
		Success: true,
		Data:    data,
	})
}
