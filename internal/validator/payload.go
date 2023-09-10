package validator

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type payloadValidator struct {
	validator *validator.Validate
	influxdb  *platform.InfluxDb
}

func NewPayloadValidator(influxdb *platform.InfluxDb) domain.PayloadValidator {
	return &payloadValidator{
		validator: validator.New(),
		influxdb:  influxdb,
	}
}

func (v *payloadValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	sid := ctx.Cookies(payload.AuthCookieKey)
	if sid == "" {
		return "", errs.New(errs.ErrAuthHeader, "missing auth header")
	}
	return sid, nil
}

func (v *payloadValidator) ValidateBody(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.BodyParser(payload); err != nil {
		return false, errs.New(errs.ErrBodyParser, err.Error())
	}
	if errors := v.validateStruct(payload); errors != nil {
		return false, errs.NewValidationErr(errs.ErrBodyValidator, "body payload is invalid", errors)
	}
	return true, nil
}

func (v *payloadValidator) ValidateQuery(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.QueryParser(payload); err != nil {
		return false, errs.New(errs.ErrQueryParser, err.Error())
	}
	if errors := v.validateStruct(payload); errors != nil {
		return false, errs.NewValidationErr(errs.ErrQueryValidator, "query payload is invalid", errors)
	}
	return true, nil
}

func (v *payloadValidator) ValidateParams(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if err := ctx.ParamsParser(payload); err != nil {
		return false, errs.New(errs.ErrParamsParser, err.Error())
	}
	if errors := v.validateStruct(payload); errors != nil {
		return false, errs.NewValidationErr(errs.ErrParamsValidator, "params payload is invalid", errors)
	}
	return true, nil
}

func (v *payloadValidator) validateStruct(payload interface{}) []errs.ValidationErrorDetail {
	var errDetails []errs.ValidationErrorDetail
	if errors := v.validator.Struct(payload); errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			detail := &errs.ValidationErrorDetail{
				Field: err.Field(),
				Type:  err.Tag(),
			}
			errDetails = append(errDetails, *detail)
		}
	}
	return errDetails
}
