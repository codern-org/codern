package domain

import (
	"io"
	"time"
)

type AuthProvider string

const (
	SelfAuth   AuthProvider = "SELF"
	GoogleAuth AuthProvider = "GOOGLE"
)

type AccountType string

const (
	FreeAccount AccountType = "FREE"
	ProAccount  AccountType = "PRO"
)

type User struct {
	Id          string       `json:"id" db:"id"`
	Email       string       `json:"email" db:"email"`
	Password    string       `json:"-" db:"password"`
	DisplayName string       `json:"displayName" db:"display_name"`
	ProfileUrl  string       `json:"profileUrl" db:"profile_url"`
	Type        AccountType  `json:"accountType" db:"account_type"`
	Provider    AuthProvider `json:"provider" db:"provider"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
}

type UpdateUser struct {
	DisplayName *string
	Profile     io.Reader
}

type UserRepository interface {
	Create(user *User) error
	Get(id string) (*User, error)
	GetBySessionId(id string) (*User, error)
	GetByEmail(email string, provider AuthProvider) (*User, error)
	Update(user *User) error
}

type UserUsecase interface {
	Create(email string, password string) (*User, error)
	CreateFromGoogle(id string, email string, name string) (*User, error)
	Get(id string) (*User, error)
	GetBySessionId(id string) (*User, error)
	GetByEmail(email string, provider AuthProvider) (*User, error)
	Update(id string, user *UpdateUser) error
	UpdatePassword(id string, oldPlainPassword string, newPlainPassword string) error
}
