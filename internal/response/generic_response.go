package response

import (
	"errors"
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Success bool           `json:"success"`
	Data    interface{}    `json:"data,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewErrParamResponse(ctx *fiber.Ctx, param string) error {
	return NewErrorResponse(
		ctx,
		fiber.StatusBadRequest,
		domain.NewError(domain.ErrParam, fmt.Sprintf("Invalid request parameter %s", param)),
	)
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error, data ...interface{}) error {
	var outputErr ErrorResponse
	var outputStatus int
	var domainError domain.DomainError

	if errors.As(err, &domainError) {
		outputErr = ErrorResponse{
			Code:    domainError.Code(),
			Message: domainError.Error(),
		}
		if data != nil {
			outputErr.Data = data
		}
		outputStatus = status
	} else {
		ctx.Locals("internal_error", err)
		outputErr = ErrorResponse{
			Code:    domain.ErrInternal,
			Message: err.Error(),
		}
		outputStatus = fiber.StatusInternalServerError
	}
	return ctx.Status(outputStatus).JSON(Response{
		Success: false,
		Error:   &outputErr,
	})
}

func NewSuccessResponse(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(Response{
		Success: true,
		Data:    data,
	})
}
