package usecase

import (
	"net/mail"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepository domain.UserRepository
	sessionUsecase domain.SessionUsecase
}

func NewUserUsecase(
	userRepository domain.UserRepository,
	sessionUsecase domain.SessionUsecase,
) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		sessionUsecase: sessionUsecase,
	}
}

func (u *userUsecase) Create(email string, password string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errs.New(
			errs.ErrInvalidEmail,
			"cannot create user with invalid email %s", email,
			err,
		)
	}

	user, err := u.GetByEmail(email, domain.SelfAuth)
	if err != nil {
		return nil, errs.New(errs.OverrideCode, "cannot create user %s", email)
	} else if user != nil {
		return nil, errs.New(errs.ErrDupEmail,
			"cannot create user due to email %s being already registered", email,
		)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return nil, errs.New(errs.ErrCreateUser, "cannot create user with invalid password", err)
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
		return nil, errs.New(errs.ErrCreateUser, "cannot create user with this email and password", err)
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
		return nil, errs.New(errs.ErrCreateUser, "cannot create user from google auth", err)
	}
	return user, nil
}

func (u *userUsecase) Get(id string) (*domain.User, error) {
	user, err := u.userRepository.Get(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetUser, "cannot get user by id %s", id, err)
	}
	return user, nil
}

func (u *userUsecase) GetBySessionId(id string) (*domain.User, error) {
	user, err := u.userRepository.GetBySessionId(id)
	if err != nil {
		return nil, errs.New(errs.ErrGetUser, "cannot get user by session id %s", id, err)
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(email string, provider domain.AuthProvider) (*domain.User, error) {
	user, err := u.userRepository.GetByEmail(email, provider)
	if err != nil {
		return nil, errs.New(errs.ErrGetUser, "cannot get user by email %s", email, err)
	}
	return user, nil
}

func (u *userUsecase) UpdatePassword(userId string, oldPlainPassword string, newPlainPassword string) error {
	user, err := u.Get(userId)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot find user id %s while update password", userId, err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPlainPassword)); err != nil {
		return errs.New(errs.ErrUserPassword, "cannot update password due to invalid old password", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPlainPassword), 10)

	if err != nil {
		return errs.New(errs.ErrUpdateUser, "cannot update password due to invalid new password", err)
	}

	user.Password = string(hashedPassword)

	err = u.userRepository.UpdatePassword(userId, user.Password)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot update password", err)
	}

	_, err = u.sessionUsecase.DestroyByUserId(userId)
	if err != nil {
		return errs.New(errs.OverrideCode, "cannot destroy session while update password", err)
	}

	return nil
}
