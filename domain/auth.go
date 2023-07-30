package domain

type AuthUsecase interface {
	Authenticate(header string) (*User, error)
	SignIn(email string, password string, ipAddress string, userAgent string) (string, error)
	SignInWithGoogle(code string, ipAddress string, userAgent string) (string, error)
	SignOut(header string) error
}
