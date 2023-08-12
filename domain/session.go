package domain

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type Session struct {
	Id        string    `json:"id" db:"id"`
	UserId    string    `json:"userId" db:"user_id"`
	IpAddress string    `json:"ipAddress" db:"ip_address"`
	UserAgent string    `json:"userAgent" db:"user_agent"`
	ExpiredAt time.Time `json:"expiredAt" db:"expired_at"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type SessionRepository interface {
	Create(session *Session) error
	Get(id string) (*Session, error)
	Delete(id string) error
	DeleteDuplicates(userId string, ipAddress string, userAgent string) error
}

type SessionUsecase interface {
	Sign(id string) string
	Unsign(header string) (string, error)
	Create(userId string, ipAddress string, userAgent string) (*fiber.Cookie, error)
	Get(header string) (*Session, error)
	Destroy(id string) error
	Validate(header string) (*Session, error)
}
