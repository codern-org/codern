package response

import (
	"errors"

	"github.com/codern-org/codern/domain"
	"github.com/gofiber/fiber/v2"
)

type GenericResponse struct {
	Success bool                 `json:"success"`
	Data    interface{}          `json:"data"`
	Error   GenericErrorResponse `json:"error"`
}

type GenericErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewErrorResponse(ctx *fiber.Ctx, status int, err error, data ...interface{}) error {
	var outputErr GenericErrorResponse
	var outputStatus int
	var genericError domain.GenericError

	if errors.As(err, &genericError) {
		outputErr = GenericErrorResponse{
			Code:    genericError.Code(),
			Message: genericError.Error(),
			Data:    data,
		}
		outputStatus = status
	} else {
		outputErr = GenericErrorResponse{
			Code:    domain.ErrInternal,
			Message: err.Error(),
		}
		outputStatus = fiber.StatusInternalServerError
	}
	return ctx.Status(outputStatus).JSON(GenericResponse{
		Success: false,
		Error:   outputErr,
	})
}

func NewSuccessResponse(ctx *fiber.Ctx, status int, data interface{}) error {
	return ctx.Status(status).JSON(GenericResponse{
		Success: true,
		Data:    data,
	})
}
