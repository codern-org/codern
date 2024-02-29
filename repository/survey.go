package repository

import (
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/platform"
)

type surveyRepository struct {
	db *platform.MySql
}

func NewSurveyRepository(db *platform.MySql) domain.SurveyRepository {
	return &surveyRepository{db: db}
}

func (r *surveyRepository) Create(survey *domain.Survey) error {
	_, err := r.db.NamedExec(
		"INSERT INTO survey (user_id, message) VALUES (:user_id, :message)",
		survey,
	)
	if err != nil {
		return fmt.Errorf("cannot create survey: %w", err)
	}
	return nil
}
