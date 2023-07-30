package usecase

import (
	"errors"

	"github.com/codern-org/codern/domain"
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

func (usecase *authUsecase) Authenticate(header string) (*domain.User, error) {
	session, err := usecase.sessionUsecase.Validate(header)
	if err != nil {
		return nil, err
	}

	user, err := usecase.userUsecase.GetBySessionId(session.Id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (usecase *authUsecase) SignIn(
	email string, password string, ipAddress string, userAgent string,
) (string, error) {
	user, err := usecase.userUsecase.GetSelfProviderUser(email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("cannot retrieve self provider user data")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("password is incorrect")
	}

	return usecase.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (usecase *authUsecase) SignInWithGoogle(
	code string, ipAddress string, userAgent string,
) (string, error) {
	token, err := usecase.googleUsecase.GetToken(code)
	if err != nil {
		return "", err
	}
	googleUser, err := usecase.googleUsecase.GetUser(token)
	if err != nil {
		return "", err
	}

	userId := usecase.userUsecase.HashId(googleUser.Id, domain.GOOGLE)
	user, err := usecase.userUsecase.Get(userId)
	if err != nil {
		return "", err
	}

	if user == nil {
		user, err = usecase.userUsecase.CreateFromGoogle(googleUser.Id, googleUser.Email)
		if err != nil {
			return "", err
		}
	}

	return usecase.sessionUsecase.Create(user.Id, ipAddress, userAgent)
}

func (usecase *authUsecase) SignOut(header string) error {
	session, err := usecase.sessionUsecase.Validate(header)
	if err != nil {
		return err
	}

	if err = usecase.sessionUsecase.Destroy(session.Id); err != nil {
		return err
	}

	return nil
}
