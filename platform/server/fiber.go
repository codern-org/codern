package server

import (
	"errors"
	"fmt"

	"github.com/codern-org/codern/controller"
	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/internal/validator"
	"github.com/codern-org/codern/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"
)

type FiberServer struct {
	app        *fiber.App
	cfg        *config.Config
	logger     *zap.Logger
	platform   *domain.Platform
	repository *domain.Repository
	usecase    *domain.Usecase
}

func NewFiberServer(
	cfg *config.Config,
	logger *zap.Logger,
	platform *domain.Platform,
	repository *domain.Repository,
	usecase *domain.Usecase,
) *FiberServer {
	return &FiberServer{
		cfg:        cfg,
		logger:     logger,
		platform:   platform,
		repository: repository,
		usecase:    usecase,
	}
}

func (s *FiberServer) Start() {
	app := fiber.New(fiber.Config{
		AppName:               s.cfg.Metadata.Name,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})
	s.app = app
	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		s.logger.Sugar().Infof("Server is listening on %s:%s", ld.Host, ld.Port)
		return nil
	})

	// Apply swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Apply middlewares
	app.Use(cors.New(cors.Config{
		AllowCredentials: constant.IsDevelopment,
		AllowOriginsFunc: func(origin string) bool {
			return constant.IsDevelopment
		},
	}))
	app.Use(requestid.New())
	app.Use(middleware.NewLogger(s.logger, s.platform.InfluxDb))

	// Apply routes
	s.applyRoutes()

	// Open fiber http server with gracefully shutdown
	if err := app.Listen(s.cfg.Client.Fiber.Address); err != nil {
		s.logger.Fatal("Server is not running", zap.Error(err))
	}
}

func (s *FiberServer) Close() error {
	return s.app.Shutdown()
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}
	return response.NewErrorResponse(ctx, code, err)
}

func (s *FiberServer) applyRoutes() {
	// Initialize Dependencies
	validator := validator.NewPayloadValidator(s.logger, s.platform.InfluxDb)

	// Initialize Middlewares
	authMiddleware := middleware.NewAuthMiddleware(s.logger, validator, s.usecase.Auth)
	workspaceMiddleware := middleware.NewWorkspaceMiddleware(s.logger, s.usecase.Workspace)

	// Initialize Controllers
	authController := controller.NewAuthController(
		s.logger, s.cfg, validator, s.usecase.Auth, s.usecase.Google, s.usecase.User,
	)
	workspaceController := controller.NewWorkspaceController(s.logger, validator, s.usecase.Workspace)

	// Initialize Routes
	api := s.app.Group("/api")

	api.Get("/auth/me", authMiddleware, authController.Me)
	api.Get("/auth/signout", authMiddleware, authController.SignOut)
	api.Post("/auth/signin", authController.SignIn)
	api.Get("/auth/google", authController.GetGoogleAuthUrl)
	api.Get("/auth/google/callback", authController.SignInWithGoogle)

	api.Get("/workspaces", authMiddleware, workspaceController.List)
	api.Get("/workspaces/:workspaceId", authMiddleware, workspaceMiddleware, workspaceController.Get)
	api.Get("/workspaces/:workspaceId/assignments", authMiddleware, workspaceMiddleware, workspaceController.ListAssignment)
	api.Get("/workspaces/:workspaceId/assignments/:assignmentId", authMiddleware, workspaceMiddleware, workspaceController.GetAssignment)
	api.Get("/workspaces/:workspaceId/assignments/:assignmentId/submissions", authMiddleware, workspaceMiddleware, workspaceController.ListSubmission)
	api.Post("/workspaces/:workspaceId/assignments/:assignmentId/submissions", authMiddleware, workspaceMiddleware, workspaceController.CreateSubmission)

	// Fallback route
	s.app.Use(func(ctx *fiber.Ctx) error {
		return response.NewErrorResponse(
			ctx,
			fiber.StatusNotFound,
			domain.NewError(domain.ErrRoute, fmt.Sprintf("No route for %s %s", ctx.Method(), ctx.Path())))
	})
}
