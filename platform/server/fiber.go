package server

import (
	"fmt"

	"github.com/codern-org/codern/domain"
	errs "github.com/codern-org/codern/domain/error"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/validator"
	"github.com/codern-org/codern/platform/server/controller"
	"github.com/codern-org/codern/platform/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/codern-org/codern/docs"
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
		ErrorHandler:          errorHandler(s.logger),
		BodyLimit:             25 * 1024 * 1024, // 25 MB
	})
	s.app = app
	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		s.logger.Sugar().Infof("Server is listening on %s:%s", ld.Host, ld.Port)
		return nil
	})

	// Apply swagger route on development mode
	if constant.IsDevelopment {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}

	// Apply middlewares
	app.Use(middleware.Cors)
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

func (s *FiberServer) applyRoutes() {
	// Initialize Dependencies
	validator := validator.NewPayloadValidator(s.platform.InfluxDb)

	// Initialize Middlewares
	fileMiddleware := middleware.NewFileMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(validator, s.usecase.Auth)
	workspaceMiddleware := middleware.NewWorkspaceMiddleware(s.usecase.Workspace)

	// Initialize Controllers
	webSocketController := controller.NewWebSocketController(s.platform.WebSocketHub)
	fileController := controller.NewFileController(s.cfg)
	authController := controller.NewAuthController(
		s.cfg, validator, s.usecase.Auth, s.usecase.Google, s.usecase.User,
	)
	workspaceController := controller.NewWorkspaceController(validator, s.usecase.Workspace)
	assignmentController := controller.NewAssignmentController(validator, s.usecase.Assignment)

	// Initialize Routes
	api := s.app.Group("/api")

	api.Get("/auth/me", authMiddleware, authController.Me)
	api.Get("/auth/signout", authMiddleware, authController.SignOut)
	api.Post("/auth/signin", authController.SignIn)
	api.Get("/auth/google", authController.GetGoogleAuthUrl)
	api.Get("/auth/google/callback", authController.SignInWithGoogle)

	api.Get("/workspaces", authMiddleware, workspaceController.List)
	api.Post("/workspaces", authMiddleware, workspaceController.CreateWorkspace)
	api.Get("/workspaces/:workspaceId", authMiddleware, workspaceMiddleware, workspaceController.Get)
	api.Post("/workspaces/:workspaceId/participants", authMiddleware, workspaceMiddleware, workspaceController.CreateParticipant)
	api.Get("/workspaces/:workspaceId/assignments", authMiddleware, workspaceMiddleware, assignmentController.List)
	api.Get("/workspaces/:workspaceId/assignments/:assignmentId", authMiddleware, workspaceMiddleware, assignmentController.Get)
	api.Get("/workspaces/:workspaceId/assignments/:assignmentId/submissions", authMiddleware, workspaceMiddleware, assignmentController.ListSubmission)
	api.Post("/workspaces/:workspaceId/assignments/:assignmentId/submissions", authMiddleware, workspaceMiddleware, assignmentController.CreateSubmission)

	// File proxy from SeaweedFS
	fs := s.app.Group("/file", authMiddleware, fileMiddleware)

	fs.Get("/user/:userId/profile", fileController.GetUserProfile)

	fs.Get("/workspaces/:workspaceId/profile", workspaceMiddleware, fileController.GetWorkspaceProfile)
	fs.Get("/workspaces/:workspaceId/assignments/:assignmentId/detail", workspaceMiddleware, fileController.GetAssignmentDetail)

	// WebSocket
	ws := s.app.Group("/ws", authMiddleware, webSocketController.Upgrade)
	ws.Get("/", webSocketController.Portal())

	// Fallback route
	s.app.Use(func(ctx *fiber.Ctx) error {
		return errs.New(errs.ErrRoute, fmt.Sprintf("No route for %s %s", ctx.Method(), ctx.Path()))
	})
}
