package server

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/response"
	"github.com/codern-org/codern/middleware"
	"github.com/codern-org/codern/platform"
	"github.com/codern-org/codern/route"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type FiberServer struct {
	cfg       *config.Config
	logger    *zap.Logger
	influxdb  *platform.InfluxDb
	mysql     *sqlx.DB
	seaweedfs *platform.SeaweedFs
	rabbitMq  *platform.RabbitMq
}

func NewFiberServer(
	cfg *config.Config,
	logger *zap.Logger,
	influxdb *platform.InfluxDb,
	mysql *sqlx.DB,
	seaweedfs *platform.SeaweedFs,
	rabbitMq *platform.RabbitMq,
) *FiberServer {
	return &FiberServer{
		cfg:       cfg,
		logger:    logger,
		influxdb:  influxdb,
		mysql:     mysql,
		seaweedfs: seaweedfs,
		rabbitMq:  rabbitMq,
	}
}

func (s *FiberServer) Start() {
	// Initialize fiber
	app := fiber.New(fiber.Config{
		AppName:               s.cfg.Metadata.Name,
		DisableStartupMessage: true,
		ErrorHandler:          errorHandler,
	})
	app.Hooks().OnListen(func(ld fiber.ListenData) error {
		s.logger.Sugar().Infof("Server is listening on %s:%s", ld.Host, ld.Port)
		return nil
	})

	route.ApplySwaggerRoutes(app)

	// Apply middlewares
	app.Use(cors.New(cors.Config{
		AllowCredentials: constant.IsDevelopment,
		AllowOriginsFunc: func(origin string) bool {
			return constant.IsDevelopment
		},
	}))
	app.Use(requestid.New())
	app.Use(middleware.NewLogger(s.logger, s.influxdb))

	// Apply routes
	route.ApplyApiRoutes(app, s.cfg, s.logger, s.influxdb, s.mysql, s.seaweedfs, s.rabbitMq)
	route.ApplyFallbackRoute(app)

	// Open fiber http server with gracefully shutdown
	go func() {
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
	s.influxdb.Close()
	s.mysql.Close()
	s.seaweedfs.Close()
	s.rabbitMq.Close()

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
