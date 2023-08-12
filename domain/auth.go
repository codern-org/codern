package domain

import "github.com/gofiber/fiber/v2"

type AuthUsecase interface {
	Authenticate(header string) (*User, error)
	SignIn(email string, password string, ipAddress string, userAgent string) (*fiber.Cookie, error)
	SignInWithGoogle(code string, ipAddress string, userAgent string) (*fiber.Cookie, error)
	SignOut(header string) error
}
