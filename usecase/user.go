package usecase

import (
	"fmt"
	"net/mail"
	"time"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/platform"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	seaweedfs      *platform.SeaweedFs
	userRepository domain.UserRepository
	sessionUsecase domain.SessionUsecase
}

func NewUserUsecase(
	seaweedfs *platform.SeaweedFs,
	userRepository domain.UserRepository,
	sessionUsecase domain.SessionUsecase,
) domain.UserUsecase {
	return &userUsecase{
		seaweedfs:      seaweedfs,
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
		return nil, errs.New(errs.SameCode, "cannot create user with email %s", email, err)
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
		return nil, errs.New(errs.ErrCreateUser, "cannot create user with email %s", email, err)
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
		return nil, errs.New(errs.ErrCreateUser, "cannot create user from google auth email %s", email, err)
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

func (u *userUsecase) Update(userId string, uu *domain.UpdateUser) error {
	user, err := u.Get(userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %d while updating user", userId, err)
	}
	if user == nil {
		return errs.New(errs.ErrGetUser, "user id %d not found", userId)
	}

	if uu.DisplayName != nil {
		user.DisplayName = *uu.DisplayName
	}
	if uu.Profile != nil {
		if user.ProfileUrl == "" {
			user.ProfileUrl = fmt.Sprintf("/user/%s/profile", user.Id)
		}
		if err := u.seaweedfs.Upload(uu.Profile, 0, user.ProfileUrl); err != nil {
			return errs.New(errs.ErrUpdateUser, "cannot upload profile of user id %s", userId, err)
		}
	}

	if err := u.userRepository.Update(user); err != nil {
		return errs.New(errs.ErrUpdateUser, "cannot update user id %s", userId, err)
	}
	return nil
}

func (u *userUsecase) UpdatePassword(userId string, oldPassword string, newPassword string) error {
	user, err := u.Get(userId)
	if err != nil {
		return errs.New(errs.SameCode, "cannot get user id %s to update password", userId, err)
	} else if user == nil {
		return errs.New(errs.ErrUserNotFound, "cannot get user id %s to update password", userId)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errs.New(errs.ErrUserPassword, "cannot update password due to invalid old password", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return errs.New(errs.ErrUpdateUser, "cannot generate new password", err)
	}
	user.Password = string(hashedPassword)

	if err := u.userRepository.Update(user); err != nil {
		return errs.New(errs.SameCode, "cannot update password", err)
	}

	if _, err = u.sessionUsecase.DestroyByUserId(userId); err != nil {
		return errs.New(errs.SameCode, "cannot destroy session while updating the password", err)
	}

	return nil
}
