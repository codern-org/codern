package usecase

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
}

func NewUserUsecase(userRepository domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepository: userRepository}
}

func (u *userUsecase) Create(email string, password string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		errMessage := fmt.Sprintf("Email %s is invalid", email)
		return nil, domain.NewError(domain.ErrInvalidEmail, errMessage)
	}

	user, err := u.GetByEmail(email, domain.SelfAuth)
	if user != nil {
		errMessage := fmt.Sprintf("Email %s is already registered", email)
		return nil, domain.NewError(domain.ErrDupEmail, errMessage)
	} else if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, err
	}

	// TODO: profile generation

	user = &domain.User{
		Id:          uuid.NewString(),
		Email:       email,
		Password:    string(hashedPassword),
		DisplayName: email,
		ProfileUrl:  "",
		Type:        domain.FreeAccount,
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
		Id:          uuid.NewString(),
		Email:       email,
		Password:    "",
		DisplayName: name,
		ProfileUrl:  "",
		Type:        domain.FreeAccount,
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
	if user == nil {
		return nil, domain.NewError(domain.ErrUserData, "Cannot get user data by this user id")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetBySessionId(id string) (*domain.User, error) {
	user, err := u.userRepository.GetBySessionId(id)
	if user == nil {
		return nil, domain.NewError(domain.ErrUserData, "Cannot get user data by this session id")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(email string, provider domain.AuthProvider) (*domain.User, error) {
	user, err := u.userRepository.GetByEmail(email, provider)
	if user == nil {
		return nil, domain.NewError(domain.ErrUserData, "Cannot get user data by this email")
	} else if err != nil {
		return nil, err
	}
	return user, nil
}
