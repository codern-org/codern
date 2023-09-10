package usecase

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"regexp"
	"strings"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type sessionUsecase struct {
	cfg               *config.Config
	sessionRepository domain.SessionRepository
}

func NewSessionUsecase(
	cfg *config.Config,
	sessionRepository domain.SessionRepository,
) domain.SessionUsecase {
	return &sessionUsecase{
		cfg:               cfg,
		sessionRepository: sessionRepository,
	}
}

func (u *sessionUsecase) Sign(id string) string {
	hmac := hmac.New(sha256.New, []byte(u.cfg.Auth.Session.Secret))
	hmac.Write([]byte(id))
	regex := regexp.MustCompile(`=+$`)
	signature := regex.ReplaceAllString(base64.StdEncoding.EncodeToString(hmac.Sum(nil)), "")
	return u.cfg.Auth.Session.Prefix + ":" + id + "." + signature
}

func (u *sessionUsecase) Unsign(header string) (string, error) {
	if !strings.HasPrefix(header, u.cfg.Auth.Session.Prefix+":") {
		return "", errs.New(errs.ErrSessionPrefix, "prefix mismatch")
	}

	id := header[len(u.cfg.Auth.Session.Prefix)+1 : strings.LastIndex(header, ".")]
	expectation := u.Sign(id)

	isLengthMatch := len([]byte(header)) == len([]byte(expectation))
	isInputMatch := subtle.ConstantTimeCompare([]byte(header), []byte(expectation)) == 1

	if !isLengthMatch || !isInputMatch {
		return "", errs.New(errs.ErrSignatureMismatch, "signature mismatch")
	}
	return id, nil
}

func (u *sessionUsecase) Create(
	userId string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	if err := u.sessionRepository.DeleteDuplicates(userId, userAgent, ipAddress); err != nil {
		return nil, err
	}

	id := uuid.NewString()
	signedId := u.Sign(id)
	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Duration(u.cfg.Auth.Session.MaxAge) * time.Second)

	err := u.sessionRepository.Create(&domain.Session{
		Id:        id,
		UserId:    userId,
		IpAddress: ipAddress,
		UserAgent: userAgent,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
	})
	if err != nil {
		return nil, err
	}

	cookie := &fiber.Cookie{
		Name:     "sid",
		Value:    signedId,
		HTTPOnly: true,
		Expires:  expiredAt,
	}
	return cookie, nil
}

func (u *sessionUsecase) Get(header string) (*domain.Session, error) {
	id, err := u.Unsign(header)
	if err != nil {
		return nil, err
	}

	session, err := u.sessionRepository.Get(id)
	if session == nil {
		return nil, errs.New(errs.ErrInvalidSession, "invalid session")
	} else if err != nil {
		return nil, err
	}
	return session, nil
}

func (u *sessionUsecase) Destroy(id string) (*fiber.Cookie, error) {
	return &fiber.Cookie{
		Name:     "sid",
		HTTPOnly: true,
		Expires:  time.Unix(0, 0),
	}, u.sessionRepository.Delete(id)
}

func (u *sessionUsecase) Validate(header string) (*domain.Session, error) {
	session, err := u.Get(header)
	if err != nil {
		return nil, err
	}

	if !time.Now().Before(session.ExpiredAt) {
		return nil, errs.New(errs.ErrSessionExpired, "session expired")
	}

	return session, nil
}
