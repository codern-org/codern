package repository

import (
	"fmt"

	"github.com/codern-org/codern/domain"
	"github.com/jmoiron/sqlx"
)

type surveyRepository struct {
	db *sqlx.DB
}

func NewSurveyRepository(db *sqlx.DB) domain.SurveyRepository {
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
