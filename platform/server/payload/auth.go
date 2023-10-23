package payload

var AuthCookieKey = "sid"

type SignInPayload struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}
