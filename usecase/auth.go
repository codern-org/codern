package usecase

import (
	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authUsecase struct {
	googleUsecase  domain.GoogleUsecase
	sessionUsecase domain.SessionUsecase
	userUsecase    domain.UserUsecase
}

func NewAuthUsecase(
	googleUsecase domain.GoogleUsecase,
	sessionUsecase domain.SessionUsecase,
	userUsecase domain.UserUsecase,
) domain.AuthUsecase {
	return &authUsecase{
		googleUsecase:  googleUsecase,
		sessionUsecase: sessionUsecase,
		userUsecase:    userUsecase,
	}
}

func (u *authUsecase) Authenticate(header string) (*domain.User, error) {
	session, err := u.sessionUsecase.Validate(header)
	if err != nil {
		return nil, err
	}
	return u.userUsecase.GetBySessionId(session.Id)
}

func (u *authUsecase) SignIn(
	email string, password string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	user, err := u.userUsecase.GetByEmail(email, domain.SelfAuth)
	if errs.HasCode(err, errs.ErrUserNotFound) {
		// Override error message
		return nil, errs.New(errs.ErrUserNotFound, "account with email %s is not registered", email)
	} else if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errs.New(errs.ErrUserPassword, "password is incorrect")
	}
	return u.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (u *authUsecase) SignInWithGoogle(
	code string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	token, err := u.googleUsecase.GetToken(code)
	if err != nil {
		return nil, err
	}
	googleUser, err := u.googleUsecase.GetUser(token)
	if err != nil {
		return nil, err
	}

	user, err := u.userUsecase.GetByEmail(googleUser.Email, domain.GoogleAuth)
	if err != nil && !errs.HasCode(err, errs.ErrUserNotFound) {
		return nil, err
	}

	if user == nil {
		user, err = u.userUsecase.CreateFromGoogle(googleUser.Id, googleUser.Email, googleUser.Name)
		if err != nil {
			return nil, err
		}
	}
	return u.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (u *authUsecase) SignOut(header string) (*fiber.Cookie, error) {
	session, err := u.sessionUsecase.Validate(header)
	if err != nil {
		return nil, err
	}
	return u.sessionUsecase.Destroy(session.Id)
}
