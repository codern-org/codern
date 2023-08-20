package controller

import (
	"net/url"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/payload"
	"github.com/codern-org/codern/internal/response"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthController struct {
	logger      *zap.Logger
	frontendCfg domain.ConfigFrontend
	validator   domain.PayloadValidator

	authUsecase   domain.AuthUsecase
	googleUsecase domain.GoogleUsecase
	userUsecase   domain.UserUsecase
}

func NewAuthController(
	logger *zap.Logger,
	frontendCfg domain.ConfigFrontend,
	validator domain.PayloadValidator,
	authUsecase domain.AuthUsecase,
	googleUsecase domain.GoogleUsecase,
	userUsecase domain.UserUsecase,
) *AuthController {
	return &AuthController{
		logger:        logger,
		frontendCfg:   frontendCfg,
		validator:     validator,
		authUsecase:   authUsecase,
		googleUsecase: googleUsecase,
		userUsecase:   userUsecase,
	}
}

// Me godoc
//
// @Summary 		Get an user data
// @Description	Get an authenticated user data
// @Tags 				auth
// @Accept 			json
// @Produce 		json
// @Security 		ApiKeyAuth
// @param 			sid header string true "Session ID"
// @Router 			/api/auth/me [get]
func (c *AuthController) Me(ctx *fiber.Ctx) error {
	return response.NewSuccessResponse(ctx, fiber.StatusOK, ctx.Locals("user"))
}

// SignIn godoc
//
// @Summary 		Sign in with self provider
// @Description Sign in with email & password provided by the user
// @Tags 				auth
// @Accept 			json
// @Produce 		json
// @Param				credentials	body	payload.AuthSignIn true "Email and password for authentication"
// @Router 			/api/auth/signin [post]
func (c *AuthController) SignIn(ctx *fiber.Ctx) error {
	var payload payload.AuthSignIn
	if ok, err := c.validator.ValidateBody(&payload, ctx); !ok {
		return err
	}

	ipAddress := ctx.IP()
	userAgent := ctx.Context().UserAgent()

	cookie, err := c.authUsecase.SignIn(payload.Email, payload.Password, ipAddress, string(userAgent))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"expired_at": cookie.Expires,
	})
}

// GetGoogleAuthUrl godoc
//
// @Summary 		Get Google auth URL
// @Description Get an url to signin with the Google account
// @Tags 				auth
// @Produce 		json
// @Router 			/api/auth/google [get]
func (c *AuthController) GetGoogleAuthUrl(ctx *fiber.Ctx) error {
	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"url": c.googleUsecase.GetOAuthUrl(),
	})
}

// SignInWithGoogle godoc
//
// @Summary 		Sign in with Google
// @Description A callback route for Google OAuth to redirect to after signing in
// @Tags 				auth
// @Produce 		json
// @Router 			/api/auth/google/callback [get]
func (c *AuthController) SignInWithGoogle(ctx *fiber.Ctx) error {
	code := ctx.Query("code")
	ipAddress := ctx.IP()
	userAgent := ctx.Context().UserAgent()

	cookie, err := c.authUsecase.SignInWithGoogle(code, ipAddress, string(userAgent))
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusBadRequest, err)
	}
	ctx.Cookie(cookie)

	url, err := url.JoinPath(c.frontendCfg.BaseUrl, c.frontendCfg.Path.SignIn)
	if err != nil {
		return err
	}

	return ctx.Redirect(url)
}

// SignOut godoc
//
// @Summary 		Sign out
// @Description Sign out and remove a sid cookie header
// @Tags 				auth
// @Produce 		json
// @Security 		ApiKeyAuths
// @param 			sid header string true "Session ID"
// @Router 			/api/auth/signout [get]
func (c *AuthController) SignOut(ctx *fiber.Ctx) error {
	sid := ctx.Cookies(payload.AuthCookieKey)

	cookie, err := c.authUsecase.SignOut(sid)
	if err != nil {
		return response.NewErrorResponse(ctx, fiber.StatusUnauthorized, err)
	}
	ctx.Cookie(cookie)

	return response.NewSuccessResponse(ctx, fiber.StatusOK, fiber.Map{
		"signout_at": time.Now(),
	})
}
