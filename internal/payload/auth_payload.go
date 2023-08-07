package payload

type AuthSignIn struct {
	Email    string `validate:"email,required"`
	Password string `validate:"required"`
}
