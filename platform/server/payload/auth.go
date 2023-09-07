package payload

var AuthCookieKey = "sid"

type SignInBody struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"required"`
}
