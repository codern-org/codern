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
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type sessionUsecase struct {
	cfgAuthSession    domain.ConfigAuthSession
	sessionRepository domain.SessionRepository
}

func NewSessionUsecase(
	cfgAuthSession domain.ConfigAuthSession,
	sessionRepository domain.SessionRepository,
) domain.SessionUsecase {
	return &sessionUsecase{
		cfgAuthSession:    cfgAuthSession,
		sessionRepository: sessionRepository,
	}
}

func (u *sessionUsecase) Sign(id string) string {
	hmac := hmac.New(sha256.New, []byte(u.cfgAuthSession.Secret))
	hmac.Write([]byte(id))
	regex := regexp.MustCompile(`=+$`)
	signature := regex.ReplaceAllString(base64.StdEncoding.EncodeToString(hmac.Sum(nil)), "")
	return u.cfgAuthSession.Prefix + ":" + id + "." + signature
}

func (u *sessionUsecase) Unsign(header string) (string, error) {
	if !strings.HasPrefix(header, u.cfgAuthSession.Prefix+":") {
		return "", domain.NewGenericError(domain.ErrSessionPrefix, "Prefix mismatch")
	}

	id := header[len(u.cfgAuthSession.Prefix)+1 : strings.LastIndex(header, ".")]
	expectation := u.Sign(id)

	isLengthMatch := len([]byte(header)) == len([]byte(expectation))
	isInputMatch := subtle.ConstantTimeCompare([]byte(header), []byte(expectation)) == 1

	if !isLengthMatch || !isInputMatch {
		return "", domain.NewGenericError(domain.ErrSignatureMismatch, "Signature mismatch")
	}
	return id, nil
}

func (u *sessionUsecase) Create(
	userId string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	err := u.sessionRepository.DeleteDuplicates(userId, userAgent, ipAddress)
	if err != nil {
		return nil, err
	}

	id := u.Sign(uuid.NewString())
	createdAt := time.Now()
	expiredAt := createdAt.Add(time.Duration(u.cfgAuthSession.MaxAge) * time.Second)

	err = u.sessionRepository.Create(&domain.Session{
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
		Value:    id,
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
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (u *sessionUsecase) Destroy(id string) error {
	return u.sessionRepository.Delete(id)
}

func (u *sessionUsecase) Validate(header string) (*domain.Session, error) {
	session, err := u.Get(header)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, domain.NewGenericError(domain.ErrInvalidSession, "Invalid session")
	}

	if !time.Now().Before(session.ExpiredAt) {
		if err := u.Destroy(session.Id); err != nil {
			return nil, err
		}
		return nil, domain.NewGenericError(domain.ErrSessionExpired, "Session expired")
	}

	return session, nil
}
