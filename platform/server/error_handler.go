package server

import (
	"errors"

	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func errorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		resStatus := fiber.StatusInternalServerError
		requestId := ctx.Locals("requestid").(string)

		var domainError *errs.DomainError
		if errors.As(err, &domainError) {
			status, ok := response.DomainErrCodeToHttpStatus[domainError.Code]
			if ok {
				resStatus = status
			}
		}

		// Mostly JSON format
		jsonErr := response.NewErrorResponse(ctx, resStatus, err)
		if jsonErr != nil {
			logger.Error(
				"Cannot marshal json response",
				zap.String("request_id", requestId),
				zap.NamedError("json_error", jsonErr),
				zap.NamedError("error", err),
			)
			// Fallback to string format if cannot marshal JSON
			return ctx.Status(500).SendString(err.Error())
		}

		return nil
	}
}
