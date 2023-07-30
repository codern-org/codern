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

func (usecase *userUsecase) HashId(id string, provider domain.AuthProvider) string {
	sha1 := sha1.New()
	sha1.Write([]byte(id + "." + string(provider)))
	return hex.EncodeToString(sha1.Sum(nil))
}

func (usecase *userUsecase) Create(email string, password string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("email %s is invalid", email)
	}

	user, err := usecase.GetSelfProviderUser(email)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return nil, fmt.Errorf("email %s already registered", email)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	// TODO: profile generation

	user = &domain.User{
		Id:          usecase.HashId(email, domain.SELF),
		Email:       email,
		Password:    string(hashedPassword),
		DisplayName: email,
		ProfilePath: "",
		Provider:    domain.SELF,
		CreatedAt:   time.Now(),
	}

	if err = usecase.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *userUsecase) CreateFromGoogle(id string, email string) (*domain.User, error) {
	// TODO: profile generation

	user := &domain.User{
		Id:          usecase.HashId(id, domain.GOOGLE),
		Email:       email,
		Password:    "",
		DisplayName: email,
		ProfilePath: "",
		Provider:    domain.SELF,
		CreatedAt:   time.Now(),
	}

	if err := usecase.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *userUsecase) Get(id string) (*domain.User, error) {
	user, err := usecase.userRepository.Get(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *userUsecase) GetBySessionId(id string) (*domain.User, error) {
	user, err := usecase.userRepository.GetBySessionId(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (usecase *userUsecase) GetSelfProviderUser(email string) (*domain.User, error) {
	user, err := usecase.userRepository.GetSelfProviderUser(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
