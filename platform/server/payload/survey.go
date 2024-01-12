package payload

type CreateSurveyPayload struct {
	Message string `json:"message" validate:"required"`
}
