package domain

import (
	"io"

	"github.com/codern-org/codern/platform"
)

type Platform struct {
	Prometheus   *platform.Prometheus
	InfluxDb     *platform.InfluxDb
	MySql        *platform.MySql
	SeaweedFs    *platform.SeaweedFs
	RabbitMq     *platform.RabbitMq
	WebSocketHub *platform.WebSocketHub
}

type Repository struct {
	Session    SessionRepository
	User       UserRepository
	Workspace  WorkspaceRepository
	Assignment AssignmentRepository
	Survey     SurveyRepository
	Misc       MiscRepository
}

type Usecase struct {
	Google     GoogleUsecase
	Session    SessionUsecase
	User       UserUsecase
	Auth       AuthUsecase
	Workspace  WorkspaceUsecase
	Assignment AssignmentUsecase
	Survey     SurveyUsecase
	Misc       MiscUsecase
}

type Publisher struct {
	Grading GradingPublisher
}

type File struct {
	Reader   io.Reader
	MimeType string
}
