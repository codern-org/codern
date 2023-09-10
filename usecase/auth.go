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
		return nil, errs.New(errs.OverrideCode, "cannot authenticate user", err)
	}
	// TODO: wrap error
	return u.userUsecase.GetBySessionId(session.Id)
}

func (u *authUsecase) SignIn(
	email string, password string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	user, err := u.userUsecase.GetByEmail(email, domain.SelfAuth)
	if err != nil {
		return nil, errs.New(errs.OverrideCode, "cannot get user data to sign in", err)
	} else if user == nil {
		return nil, errs.New(errs.ErrUserNotFound, "account with email %s is not registered", email)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errs.New(errs.ErrUserPassword, "password is incorrect", err)
	}
	// TODO: wrap error
	return u.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (u *authUsecase) SignInWithGoogle(
	code string, ipAddress string, userAgent string,
) (*fiber.Cookie, error) {
	token, err := u.googleUsecase.GetToken(code)
	if err != nil {
		return nil, errs.New(errs.OverrideCode, "cannot sign in with google", err)
	}
	googleUser, err := u.googleUsecase.GetUser(token)
	if err != nil {
		return nil, errs.New(errs.OverrideCode, "cannot sign in with google", err)
	}

	user, err := u.userUsecase.GetByEmail(googleUser.Email, domain.GoogleAuth)
	if err != nil && !errs.HasCode(err, errs.ErrUserNotFound) {
		return nil, errs.New(errs.OverrideCode, "cannot get user data to sign in with google", err)
	}

	if user == nil {
		user, err = u.userUsecase.CreateFromGoogle(googleUser.Id, googleUser.Email, googleUser.Name)
		if err != nil {
			return nil, errs.New(errs.OverrideCode, "cannot create user to sign in with google", err)
		}
	}
	// TODO: wrap error
	return u.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (u *authUsecase) SignOut(header string) (*fiber.Cookie, error) {
	session, err := u.sessionUsecase.Validate(header)
	if err != nil {
		return nil, errs.New(errs.OverrideCode, "cannot validate session to signout", err)
	}
	// TODO: wrap error
	return u.sessionUsecase.Destroy(session.Id)
}
