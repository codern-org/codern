package route

import (
	"github.com/codern-org/codern/controller"
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/validator"
	"github.com/codern-org/codern/middleware"
	"github.com/codern-org/codern/repository"
	"github.com/codern-org/codern/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func ApplyApiRoutes(
	app *fiber.App,
	cfg *domain.Config,
	logger *zap.Logger,
	influxdb domain.InfluxDb,
	mysql *sqlx.DB,
) {
	// Initialize Dependencies
	validator := validator.NewPayloadValidator(logger, influxdb)

	// Initialize Repositories
	sessionRepository := repository.NewSessionRepository(mysql)
	userRepository := repository.NewUserRepository(mysql)

	// Initialize Usercases
	googleUsecase := usecase.NewGoogleUsecase(cfg.Google)
	sessionUsecase := usecase.NewSessionUsecase(cfg.Auth.Session, sessionRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)
	authUsecase := usecase.NewAuthUsecase(googleUsecase, sessionUsecase, userUsecase)

	// Initialize Controllers
	authController := controller.NewAuthController(
		logger, cfg.Client.Frontend, validator, authUsecase, googleUsecase, userUsecase,
	)

	// Initialize Middlewares
	authMiddleware := middleware.NewAuthMiddleware(logger, validator, authUsecase)

	// Initialize Routes
	auth := app.Group("/api/auth")
	auth.Get("/me", authMiddleware, authController.Me)
	auth.Post("/signin", authController.SignIn)
	auth.Get("/signout", authMiddleware, authController.SignOut)

	auth.Get("/google", authController.GetGoogleAuthUrl)
	auth.Get("/google/callback", authController.SignInWithGoogle)
}
