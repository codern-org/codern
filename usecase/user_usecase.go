package usecase

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/mail"
	"time"

	"github.com/codern-org/codern/domain"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) HashId(id string, provider domain.AuthProvider) string {
	sha1 := sha1.New()
	sha1.Write([]byte(id + "." + string(provider)))
	return hex.EncodeToString(sha1.Sum(nil))
}

func (u *userUsecase) Create(email string, password string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		errMessage := fmt.Sprintf("email %s is invalid", email)
		return nil, domain.NewGenericError(domain.ErrInvalidEmail, errMessage)
	}

	user, err := u.GetSelfProviderUser(email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		errMessage := fmt.Sprintf("email %s is already registered", email)
		return nil, domain.NewGenericError(domain.ErrDupEmail, errMessage)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	// TODO: profile generation

	user = &domain.User{
		Id:          u.HashId(email, domain.SelfAuth),
		Email:       email,
		Password:    string(hashedPassword),
		DisplayName: "",
		ProfileUrl:  "",
		Provider:    domain.SelfAuth,
		CreatedAt:   time.Now(),
	}

	if err = u.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) CreateFromGoogle(id string, email string, name string) (*domain.User, error) {
	// TODO: profile generation

	user := &domain.User{
		Id:          u.HashId(id, domain.GoogleAuth),
		Email:       email,
		Password:    "",
		DisplayName: name,
		ProfileUrl:  "",
		Provider:    domain.GoogleAuth,
		CreatedAt:   time.Now(),
	}

	if err := u.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) Get(id string) (*domain.User, error) {
	user, err := u.userRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetBySessionId(id string) (*domain.User, error) {
	user, err := u.userRepository.GetBySessionId(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetSelfProviderUser(email string) (*domain.User, error) {
	user, err := u.userRepository.GetSelfProviderUser(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
