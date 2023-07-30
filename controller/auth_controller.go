package controller

import (
	"github.com/codern-org/codern/domain"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthController struct {
	logger *zap.Logger

	authUsecase   domain.AuthUsecase
	googleUsecase domain.GoogleUsecase
	userUsecase   domain.UserUsecase
}

func NewAuthContoller(
	logger *zap.Logger,
	authUsecase domain.AuthUsecase,
	googleUsecase domain.GoogleUsecase,
	userUsecase domain.UserUsecase,
) *AuthController {
	return &AuthController{
		logger:        logger,
		authUsecase:   authUsecase,
		googleUsecase: googleUsecase,
		userUsecase:   userUsecase,
	}
}

func (c *AuthController) Me(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "me",
	})
}

func (c *AuthController) SignIn(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "signin",
	})
}

func (c *AuthController) SignOut(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "signout",
	})
}

func (c *AuthController) GetGoogleAuthUrl(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"url": c.googleUsecase.GetOAuthUrl(),
	})
}

func (c *AuthController) SignInWithGoogle(ctx *fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"message": "google",
	})
}
