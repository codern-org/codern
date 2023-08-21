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
	workspaceRepository := repository.NewWorkspaceRepository(mysql)

	// Initialize Usercases
	googleUsecase := usecase.NewGoogleUsecase(cfg.Google)
	sessionUsecase := usecase.NewSessionUsecase(cfg.Auth.Session, sessionRepository)
	userUsecase := usecase.NewUserUsecase(userRepository)
	authUsecase := usecase.NewAuthUsecase(googleUsecase, sessionUsecase, userUsecase)
	workspaceUsecase := usecase.NewWorkspaceUsecase(workspaceRepository)

	// Initialize Controllers
	authController := controller.NewAuthController(
		logger, cfg.Client.Frontend, validator, authUsecase, googleUsecase, userUsecase,
	)
	workspaceController := controller.NewWorkspaceController(logger, workspaceUsecase)

	// Initialize Middlewares
	authMiddleware := middleware.NewAuthMiddleware(logger, validator, authUsecase)

	// Initialize Routes
	api := app.Group("/api")

	auth := api.Group("/auth")
	workspace := api.Group("/workspace", authMiddleware)

	auth.Get("/me", authMiddleware, authController.Me)
	auth.Get("/signout", authMiddleware, authController.SignOut)
	auth.Post("/signin", authController.SignIn)
	auth.Get("/google", authController.GetGoogleAuthUrl)
	auth.Get("/google/callback", authController.SignInWithGoogle)

	workspace.Get("/", workspaceController.GetAllFromUserId)
	workspace.Get("/:id", workspaceController.Get)
}
