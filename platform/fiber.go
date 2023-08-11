package platform

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/middleware"
	"github.com/codern-org/codern/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type fiberServer struct {
	cfg      *domain.Config
	logger   *zap.Logger
	influxdb domain.InfluxDb
	mysql    *sqlx.DB
}

func NewFiberServer(
	cfg *domain.Config,
	logger *zap.Logger,
	influxdb domain.InfluxDb,
	mysql *sqlx.DB,
) domain.FiberServer {
	return &fiberServer{
		cfg:      cfg,
		logger:   logger,
		influxdb: influxdb,
		mysql:    mysql,
	}
}

func (s *fiberServer) Start() {
	// Initialize fiber
	app := fiber.New(fiber.Config{
		AppName:               s.cfg.Metadata.Name,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})

	route.ApplySwaggerRoutes(app)

	// Apply middlewares
	app.Use(requestid.New())
	app.Use(middleware.NewLogger(s.logger, s.influxdb))

	// Apply routes
	route.ApplyApiRoutes(app, s.cfg, s.logger, s.influxdb, s.mysql)
	route.ApplyFallbackRoute(app)

	// Open fiber http server with gracefully shutdown
	go func() {
		s.logger.Info("Server is starting")
		if err := app.Listen(s.cfg.Client.Fiber.Address); err != nil {
			s.logger.Panic("Server is not running", zap.Error(err))
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	// Block the main thread until an interrupt is received
	<-signals
	s.logger.Info("Server is shutting down")
	if err := app.Shutdown(); err != nil {
		s.logger.Error("Server is not shutting down", zap.Error(err))
	}
	s.logger.Info("Running cleanup tasks")

	// Clean up

	s.logger.Info("Server was successful shutdown")
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	return response.NewErrorResponse(ctx, code, err)
}
