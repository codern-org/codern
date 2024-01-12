package usecase

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
)

type surveyUsecase struct {
	surveyRepository domain.SurveyRepository
}

func NewSurveyUsecase(
	surveyRepository domain.SurveyRepository,
) domain.SurveyUsecase {
	return &surveyUsecase{
		surveyRepository: surveyRepository,
	}
}

func (u *surveyUsecase) Create(userId string, message string) error {
	survey := &domain.Survey{
		UserId:  userId,
		Message: message,
	}

	err := u.surveyRepository.Create(survey)
	if err != nil {
		return errs.New(errs.ErrCreateSurvey, "cannot create survey")
	}

	return nil
}
