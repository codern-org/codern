package domain

import "time"

type Survey struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	Message   string    `json:"message" db:"message"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type SurveyRepository interface {
	Create(survey *Survey) error
}

type SurveyUsecase interface {
	Create(userId string, message string) error
}
