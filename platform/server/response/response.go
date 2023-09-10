package response

import (
	"errors"

	errs "github.com/codern-org/codern/domain/error"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error) error {
	var domainError *errs.DomainError
	if errors.As(err, &domainError) {
		return ctx.Status(status).JSON(Response{
			Success: false,
			Error:   domainError,
		})
	}

	return ctx.Status(status).JSON(Response{
		Success: false,
		Error: &errs.DomainError{
			Code:    errs.ErrInternal,
			Message: "Internal server error",
		},
	})
}

func NewSuccessResponse(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(Response{
		Success: true,
		Data:    data,
	})
}
