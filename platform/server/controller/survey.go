package controller

import (
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/codern-org/codern/platform/server/payload"
	"github.com/codern-org/codern/platform/server/response"
	"github.com/gofiber/fiber/v2"
)

type SurveyController struct {
	validator     domain.PayloadValidator
	surveyUsecase domain.SurveyUsecase
}

func NewSurveyController(
	validator domain.PayloadValidator,
	SurveyUsecase domain.SurveyUsecase,
) *SurveyController {
	return &SurveyController{
		validator:     validator,
		surveyUsecase: SurveyUsecase,
	}
}

func (c *SurveyController) CreateSurvey(ctx *fiber.Ctx) error {
	var payload payload.CreateSurveyPayload
	if ok, err := c.validator.Validate(&payload, ctx); !ok {
		return err
	}

	user := middleware.GetUserFromCtx(ctx)

	err := c.surveyUsecase.Create(user.Id, payload.Message)
	if err != nil {
		return err
	}

	return response.NewSuccessResponse(ctx, fiber.StatusCreated, fiber.Map{})
}
