package validator

import (
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/payload"
	"github.com/codern-org/codern/internal/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type payloadValidator struct {
	logger    *zap.Logger
	influxdb  domain.InfluxDb
	validator *validator.Validate
}

func NewPayloadValidator(
	logger *zap.Logger,
	influxdb domain.InfluxDb,
) domain.PayloadValidator {
	return &payloadValidator{
		logger:    logger,
		influxdb:  influxdb,
		validator: validator.New(),
	}
}

func (v payloadValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	sid := ctx.Get(payload.AuthCookieKey)
	if sid == "" {
		v.logger.Warn(
			"Missing auth header",
			zap.String("path", ctx.Path()),
			zap.String("ip_address", ctx.IP()),
			zap.String("user_agent", string(ctx.Context().UserAgent())),
		)
		return "", ctx.
			Status(fiber.StatusBadRequest).
			JSON(response.GenericErrorResponse{
				Code:    response.ErrAuthHeaderNotFound,
				Message: fmt.Sprintf("Missing `%s` header for authentication", payload.AuthCookieKey),
			})
	}
	return sid, nil
}

func (v payloadValidator) ValidateBody(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.BodyParser(&payload); err != nil {
		return false, ctx.
			Status(fiber.StatusUnprocessableEntity).
			JSON(response.GenericResponse{
				Sucess: false,
				Errors: []response.GenericErrorResponse{
					{
						Code:    response.ErrPayloadParser,
						Message: err.Error(),
					},
				},
			})
	}

	if errs := v.validateStruct(payload); errs != nil {
		v.logger.Warn("Payload validation failed", zap.Any("details", errs))
		return false, ctx.
			Status(fiber.StatusBadRequest).
			JSON(response.GenericResponse{
				Sucess: false,
				Errors: []response.GenericErrorResponse{
					{
						Code:    response.ErrPayloadValidator,
						Message: "Payload validation failed",
						Details: errs,
					},
				},
			})
	}

	return true, nil
}

func (v payloadValidator) validateStruct(payload interface{}) []response.GenericValidationError {
	errors := []response.GenericValidationError{}

	if errs := v.validator.Struct(payload); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			var element response.GenericValidationError
			element.Namespace = err.StructNamespace()
			element.Field = err.Tag()
			element.Value = err.Value()
			errors = append(errors, element)
		}
	}

	return errors
}
