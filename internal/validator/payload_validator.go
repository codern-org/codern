package validator

import (
	"log"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/payload"
	"github.com/codern-org/codern/internal/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type payloadValidator struct {
	logger    *zap.Logger
	validator *validator.Validate
	influxdb  domain.InfluxDb
}

func NewPayloadValidator(
	logger *zap.Logger,
	influxdb domain.InfluxDb,
) domain.PayloadValidator {
	return &payloadValidator{
		logger:    logger,
		validator: validator.New(),
		influxdb:  influxdb,
	}
}

func (v *payloadValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	sid := ctx.Cookies(payload.AuthCookieKey)
	if sid == "" {
		v.logger.Warn(
			"Missing auth header",
			zap.String("path", ctx.Path()),
			zap.String("ip_address", ctx.IP()),
			zap.String("user_agent", string(ctx.Context().UserAgent())),
		)
		return "", response.NewErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			domain.NewError(domain.ErrAuthHeader, "missing auth header"),
		)
	}
	return sid, nil
}

func (v *payloadValidator) ValidateBody(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.BodyParser(&payload); err != nil {
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusUnprocessableEntity,
			domain.NewError(domain.ErrBodyParser, err.Error()),
		)
	}

	if errs := v.validateStruct(payload); errs != nil {
		v.logger.Warn("Body payload validation failed", zap.Any("details", errs))
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			domain.NewErrorWithData(domain.ErrBodyValidator, "body payload is invalid", errs),
		)
	}

	return true, nil
}

func (v *payloadValidator) ValidateQuery(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.QueryParser(payload); err != nil {
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusUnprocessableEntity,
			domain.NewError(domain.ErrQueryParser, err.Error()),
		)
	}

	if errs := v.validateStruct(payload); errs != nil {
		v.logger.Warn("Query payload validation failed", zap.Any("details", errs))
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			domain.NewErrorWithData(domain.ErrQueryValidator, "query payload is invalid", errs),
		)
	}

	return true, nil
}

func (v *payloadValidator) ValidateParams(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.ParamsParser(payload); err != nil {
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusUnprocessableEntity,
			domain.NewError(domain.ErrParamsParser, err.Error()),
		)
	}

	if errs := v.validateStruct(payload); errs != nil {
		v.logger.Warn("Params payload validation failed", zap.Any("details", errs))
		return false, response.NewErrorResponse(
			ctx,
			fiber.StatusBadRequest,
			domain.NewErrorWithData(domain.ErrParamsValidator, "params payload is invalid", errs),
		)
	}

	return true, nil
}

func (v *payloadValidator) validateStruct(payload interface{}) []*domain.ValidationError {
	var errors []*domain.ValidationError

	if errs := v.validator.Struct(payload); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErr := domain.NewValidationError(
				err.StructNamespace(),
				err.Tag(),
				err.Value(),
			)
			log.Printf("%+v", validationErr)
			errors = append(errors, validationErr)
		}
	}

	return errors
}
