package usecase

import (
	"github.com/codern-org/codern/domain"
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
	user, err := u.userUsecase.GetSelfProviderUser(email)
	if domain.HasErrorCode(err, domain.ErrUserData) {
		// Override error message
		return nil, domain.NewError(domain.ErrUserData, "This account is not registered")
	} else if err != nil {
		return nil, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, domain.NewError(domain.ErrUserPassword, "Password is incorrect")
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

	userId := u.userUsecase.HashId(googleUser.Id, domain.GoogleAuth)
	user, err := u.userUsecase.Get(userId)
	if err != nil && !domain.HasErrorCode(err, domain.ErrUserData) {
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
