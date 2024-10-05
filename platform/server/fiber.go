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
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/swagger"
	"go.uber.org/zap"

	_ "github.com/codern-org/codern/other/swagger"
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

		BodyLimit:    100 * 1024 * 1024, // 100 MB
		ServerHeader: "codern",

		TrustedProxies:          s.cfg.Client.Fiber.TrustedProxies,
		EnableTrustedProxyCheck: len(s.cfg.Client.Fiber.TrustedProxies) > 0,
		ProxyHeader:             s.cfg.Client.Fiber.ProxyHeader,
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
	app.Use(favicon.New())
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
	validator := validator.NewPayloadValidator()

	// Initialize Middlewares
	fileMiddleware := middleware.NewFileMiddleware()
	authMiddleware := middleware.NewAuthMiddleware(validator, s.usecase.Auth)
	publishableWorkspaceMiddleware := middleware.NewPublishableWorkspaceMiddleware(validator, s.usecase.Auth, s.usecase.Workspace)
	workspaceMiddleware := middleware.NewWorkspaceMiddleware(validator, s.usecase.Workspace)
	scoreboardMiddleware := middleware.NewScoreboardMiddleware(validator, s.usecase.Auth, s.usecase.Workspace, s.usecase.Misc)

	// Initialize Controllers
	healtController := controller.NewHealthController(s.cfg)
	webSocketController := controller.NewWebSocketController(s.platform.WebSocketHub)
	fileController := controller.NewFileController(s.cfg, validator, s.usecase.Workspace)
	authController := controller.NewAuthController(
		s.cfg, validator, s.usecase.Auth, s.usecase.Google, s.usecase.User,
	)
	workspaceController := controller.NewWorkspaceController(validator, s.usecase.Workspace)
	assignmentController := controller.NewAssignmentController(validator, s.usecase.Assignment)
	userController := controller.NewUserController(validator, s.usecase.User)
	surveyController := controller.NewSurveyController(validator, s.usecase.Survey)

	// Initialize Routes
	api := s.app.Group("/")

	api.Get("/", middleware.PathType("healthcheck"), healtController.Index)
	api.Get("/health", middleware.PathType("healthcheck"), healtController.Check)
	api.Get("/metrics", middleware.PathType("healthcheck"), healtController.Metrics)

	auth := api.Group("/auth", middleware.PathType("auth"))
	auth.Get("/me", authMiddleware, authController.Me)
	auth.Get("/signout", authMiddleware, authController.SignOut)
	auth.Post("/signin", authController.SignIn)
	auth.Get("/google", authController.GetGoogleAuthUrl)
	auth.Get("/google/callback", authController.SignInWithGoogle)

	user := api.Group("/users", middleware.PathType("user"))
	user.Patch("/", authMiddleware, userController.Update)
	user.Patch("/password", authMiddleware, userController.UpdatePassword)

	workspace := api.Group("/workspaces", middleware.PathType("workspace"))
	workspace.Get("/join/:invitationId", authMiddleware, workspaceController.JoinByInvitationCode)
	workspace.Get("/", authMiddleware, workspaceMiddleware, workspaceController.List)
	workspace.Post("/", authMiddleware, workspaceMiddleware, workspaceController.Create)
	workspace.Patch("/:workspaceId", authMiddleware, workspaceMiddleware, workspaceController.Update)
	workspace.Delete("/:workspaceId", authMiddleware, workspaceMiddleware, workspaceController.Delete)
	workspace.Get("/:workspaceId", publishableWorkspaceMiddleware, workspaceController.Get)
	workspace.Get("/:workspaceId/participants", authMiddleware, workspaceMiddleware, workspaceController.ListParticipant)
	workspace.Patch("/:workspaceId/participants/:userId", authMiddleware, workspaceMiddleware, workspaceController.UpdateParticipant)
	workspace.Delete("/:workspaceId/participants/:userId", authMiddleware, workspaceMiddleware, workspaceController.DeleteParticipant)
	workspace.Get("/:workspaceId/scoreboard", scoreboardMiddleware, cache.New(), workspaceController.GetScoreboard)

	assignment := workspace.Group("/:workspaceId/assignments")
	assignment.Get("/", authMiddleware, workspaceMiddleware, assignmentController.List)
	assignment.Post("/", authMiddleware, workspaceMiddleware, assignmentController.Create)
	assignment.Get("/:assignmentId", authMiddleware, workspaceMiddleware, assignmentController.Get)
	assignment.Patch("/:assignmentId", authMiddleware, workspaceMiddleware, assignmentController.Update)
	assignment.Delete("/:assignmentId", authMiddleware, workspaceMiddleware, assignmentController.Delete)
	assignment.Get("/:assignmentId/submissions", authMiddleware, workspaceMiddleware, assignmentController.ListSubmission)
	assignment.Post("/:assignmentId/submissions", authMiddleware, workspaceMiddleware, assignmentController.CreateSubmission)

	invitation := workspace.Group("/:workspaceId/invitation", middleware.PathType("invitation"))
	invitation.Get("/", authMiddleware, workspaceMiddleware, workspaceController.GetInvitations)
	invitation.Post("/", authMiddleware, workspaceMiddleware, workspaceController.CreateInvitation)
	invitation.Delete("/:invitationId", authMiddleware, workspaceMiddleware, workspaceController.DeleteInvitation)

	survey := s.app.Group("/survey")
	survey.Post("/", authMiddleware, surveyController.CreateSurvey)

	// File proxy from SeaweedFS
	fs := s.app.Group("/file", middleware.PathType("file"), fileMiddleware)
	fs.Get("/user/:userId/profile", fileController.GetUserProfile)
	fs.Get("/workspaces/:workspaceId/profile", fileController.GetWorkspaceProfile)
	fs.Get("/workspaces/:workspaceId/assignments/:assignmentId/detail/*", authMiddleware, workspaceMiddleware, fileController.GetAssignmentDetail)
	fs.Get("/workspaces/:workspaceId/assignments/:assignmentId/testcase/:testcaseFile", authMiddleware, workspaceMiddleware, fileController.GetAssignmentTestcase)
	fs.Get("/workspaces/:workspaceId/assignments/:assignmentId/submissions/:userId/:submissionId", authMiddleware, workspaceMiddleware, fileController.GetSubmission)

	// WebSocket
	ws := s.app.Group("/ws", authMiddleware, webSocketController.Upgrade)
	ws.Get("/", webSocketController.Portal())

	// Fallback route
	s.app.Use(func(ctx *fiber.Ctx) error {
		return errs.New(errs.ErrRoute, fmt.Sprintf("No route for %s %s", ctx.Method(), ctx.Path()))
	})
}
