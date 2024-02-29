package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/constant"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type payloadValidator struct {
	validator *validator.Validate
}

func NewPayloadValidator() domain.PayloadValidator {
	return &payloadValidator{
		validator: validator.New(),
	}
}

func (v *payloadValidator) ValidateAuth(ctx *fiber.Ctx) (string, error) {
	sid := ctx.Cookies(constant.SessionCookieName)
	if sid == "" {
		return "", errs.New(errs.ErrAuthHeader, "missing auth header")
	}
	return sid, nil
}

func (v *payloadValidator) Validate(payload interface{}, ctx *fiber.Ctx) (bool, error) {
	if len(ctx.Body()) != 0 { // Prevent ErrUnprocessableEntity when body is empty
		if err := ctx.BodyParser(payload); err != nil {
			return false, errs.New(errs.ErrBodyParser, err.Error())
		}
	}
	if err := ctx.ParamsParser(payload); err != nil {
		return false, errs.New(errs.ErrParamsParser, err.Error())
	}
	if err := ctx.QueryParser(payload); err != nil {
		return false, errs.New(errs.ErrQueryParser, err.Error())
	}
	if err := fileParser(payload, ctx); err != nil {
		return false, errs.New(errs.ErrBodyParser, err.Error())
	}

	if errors := v.validateStruct(payload); errors != nil {
		return false, errs.NewPayloadError(errors)
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

func fileParser(payload interface{}, ctx *fiber.Ctx) error {
	v := reflect.ValueOf(payload)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("interface must be a pointer to struct")
	}
	v = v.Elem() // Unwrap interfae or pointer

	var form *multipart.Form
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		fileKey := field.Tag.Get("file")
		if fileKey != "" {
			// Parse multipart form if not parsed before
			if form == nil {
				ctxForm, err := ctx.MultipartForm()
				if err != nil {
					return fmt.Errorf("cannot parse the file")
				}
				form = ctxForm
			}

			if field.Type == reflect.TypeOf((*multipart.File)(nil)).Elem() {
				// Parse a single file.
				// If the payload contains multiple files, the first file is being parsed
				for name, headers := range form.File {
					if name == fileKey {
						file, err := headers[0].Open()
						if err != nil {
							return fmt.Errorf("cannot parse the file")
						}
						v.Field(i).Set(reflect.ValueOf(file))
					}
				}
			} else if field.Type == reflect.TypeOf((*[]multipart.File)(nil)).Elem() {
				// Parse multiple files
				for name, headers := range form.File {
					if name != fileKey {
						continue
					}
					var files []multipart.File
					for _, header := range headers {
						file, err := header.Open()
						if err != nil {
							return fmt.Errorf("cannot parse the file")
						}
						files = append(files, file)
					}
					v.Field(i).Set(reflect.ValueOf(files))
					break
				}
			}
		}
	}

	return nil
}
