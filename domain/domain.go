package domain

import (
	"github.com/codern-org/codern/platform"
	"github.com/jmoiron/sqlx"
)

type Platform struct {
	InfluxDb     *platform.InfluxDb
	MySql        *sqlx.DB
	SeaweedFs    *platform.SeaweedFs
	RabbitMq     *platform.RabbitMq
	WebSocketHub *platform.WebSocketHub
}

type Repository struct {
	Session   SessionRepository
	User      UserRepository
	Workspace WorkspaceRepository
}

type Usecase struct {
	Google    GoogleUsecase
	Session   SessionUsecase
	User      UserUsecase
	Auth      AuthUsecase
	Workspace WorkspaceUsecase
}

type Publisher struct {
	Grading GradingPublisher
}
