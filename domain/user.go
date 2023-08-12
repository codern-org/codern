package domain

import "time"

type AuthProvider string

const (
	SELF   AuthProvider = "SELF"
	GOOGLE AuthProvider = "GOOGLE"
)

type User struct {
	Id          string       `json:"id" db:"id"`
	Email       string       `json:"email" db:"email"`
	Password    string       `json:"-" db:"password"`
	DisplayName string       `json:"displayName" db:"display_name"`
	ProfileUrl  string       `json:"profileUrl" db:"profile_url"`
	Provider    AuthProvider `json:"provider" db:"provider"`
	CreatedAt   time.Time    `json:"createdAt" db:"created_at"`
}

type UserRepository interface {
	Create(user *User) error
	Get(id string) (*User, error)
	GetBySessionId(id string) (*User, error)
	GetSelfProviderUser(email string) (*User, error)
}

type UserUsecase interface {
	HashId(id string, provider AuthProvider) string
	Create(email string, password string) (*User, error)
	CreateFromGoogle(id string, email string, name string) (*User, error)
	Get(id string) (*User, error)
	GetBySessionId(id string) (*User, error)
	GetSelfProviderUser(email string) (*User, error)
}
