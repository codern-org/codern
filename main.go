package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codern-org/codern/domain"
	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/logger"
	"github.com/codern-org/codern/platform"
	"github.com/codern-org/codern/platform/amqp/consumer"
	"github.com/codern-org/codern/platform/amqp/publisher"
	"github.com/codern-org/codern/platform/server"
	"github.com/codern-org/codern/repository"
	"github.com/codern-org/codern/usecase"
	"go.uber.org/zap"
)

// @title Codern API Server
// @version 0.0.0
// @description The API Server of Codern
//
// @securityDefinitions.apikey 	ApiKeyAuth
// @in 													header
// @name												Authorization
func main() {
	// Initialize logger
	logger := logger.NewLogger()
	defer logger.Sync()

	logger.Sugar().Infof("Codern API server %s", constant.Version)

	if constant.IsDevelopment {
		logger.Warn("Running in development mode")
	}

	// Load flags
	var configPath string
	flag.StringVar(&configPath, "config", "./config/config.yaml", "path to a config file")
	flag.Parse()

	// Load configuration file
	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal("Cannot load a config file", zap.Error(err))
	}
	logger.Info("Configuration file loaded successfully")

	// Initialize dependencies
	platform := initPlatform(cfg, logger)
	repository := initRepository(platform.MySql)
	publisher := initPublisher(cfg, platform)
	usecase := initUsecase(cfg, logger, platform, repository, publisher)

	startConsumer(logger, platform, usecase)

	// Initialize server with gracefully shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	fiber := server.NewFiberServer(cfg, logger, platform, repository, usecase)
	go fiber.Start()

	// Block the main thread until an interrupt is received
	<-signals
	logger.Info("Server is shutting down")
	if err := fiber.Close(); err != nil {
		logger.Error("Server is not shutting down", zap.Error(err))
	}
	logger.Info("Running cleanup tasks")

	// Clean up
	platform.RabbitMq.Close()
	platform.SeaweedFs.Close()
	platform.MySql.Close()
	platform.InfluxDb.Close()

	logger.Info("Server was successful shutdown")
}

func initPlatform(cfg *config.Config, logger *zap.Logger) *domain.Platform {
	prometheus := platform.NewPrometheus()

	start := time.Now()
	influxdb, err := platform.NewInfluxDb(
		cfg.Client.InfluxDb.Url,
		cfg.Client.InfluxDb.Token,
		cfg.Client.InfluxDb.Org,
		cfg.Client.InfluxDb.Bucket,
		logger,
	)
	if err != nil {
		logger.Error("Cannot create a InfluxDB connection", zap.Error(err))
	} else {
		logger.Info("Connected to InfluxDB", zap.String("connection_time", time.Since(start).String()))
	}

	start = time.Now()
	mysql, err := platform.NewMySql(cfg.Client.MySql.Uri)
	if err != nil {
		logger.Fatal("Cannot open MySQL database connection", zap.Error(err))
	}
	logger.Info("Connected to MySQL", zap.String("connection_time", time.Since(start).String()))

	start = time.Now()
	seaweedfs, err := platform.NewSeaweedFs(
		cfg.Client.SeaweedFs.MasterUrl,
		cfg.Client.SeaweedFs.FilerUrls.Internal,
	)
	if err != nil {
		logger.Fatal("Cannot open SeaweedFs connection", zap.Error(err))
	}
	logger.Info("Connected to SeaweedFs", zap.String("connection_time", time.Since(start).String()))

	start = time.Now()
	rabbitmq, err := platform.NewRabbitMq(cfg.Client.RabbitMq.Url)
	if err != nil {
		logger.Fatal("Cannot open RabbitMq connection", zap.Error(err))
	}
	logger.Info("Connected to RabbitMq", zap.String("connection_time", time.Since(start).String()))

	webSocketHub := platform.NewWebSocketHub(prometheus)

	return &domain.Platform{
		Prometheus:   prometheus,
		InfluxDb:     influxdb,
		MySql:        mysql,
		SeaweedFs:    seaweedfs,
		RabbitMq:     rabbitmq,
		WebSocketHub: webSocketHub,
	}
}

func initRepository(mysql *platform.MySql) *domain.Repository {
	return &domain.Repository{
		Session:    repository.NewSessionRepository(mysql),
		User:       repository.NewUserRepository(mysql),
		Workspace:  repository.NewWorkspaceRepository(mysql),
		Assignment: repository.NewAssignmentRepository(mysql),
		Survey:     repository.NewSurveyRepository(mysql),
		Misc:       repository.NewMiscRepsitory(mysql),
	}
}

func initUsecase(
	cfg *config.Config,
	logger *zap.Logger,
	platform *domain.Platform,
	repository *domain.Repository,
	publisher *domain.Publisher,
) *domain.Usecase {
	miscUsecase := usecase.NewMiscUsecase(repository.Misc)
	googleUsecase := usecase.NewGoogleUsecase(cfg)
	sessionUsecase := usecase.NewSessionUsecase(cfg, repository.Session)
	userUsecase := usecase.NewUserUsecase(platform.SeaweedFs, repository.User, sessionUsecase)
	authUsecase := usecase.NewAuthUsecase(googleUsecase, sessionUsecase, userUsecase)
	workspaceUsecase := usecase.NewWorkspaceUsecase(platform.SeaweedFs, repository.Workspace, repository.User, userUsecase)
	assignmentUsecase := usecase.NewAssignmentUsecase(platform.SeaweedFs, repository.Assignment, publisher.Grading, workspaceUsecase)
	surveyUsecase := usecase.NewSurveyUsecase(repository.Survey)

	return &domain.Usecase{
		Google:     googleUsecase,
		Session:    sessionUsecase,
		User:       userUsecase,
		Auth:       authUsecase,
		Workspace:  workspaceUsecase,
		Assignment: assignmentUsecase,
		Survey:     surveyUsecase,
		Misc:       miscUsecase,
	}
}

func initPublisher(
	cfg *config.Config,
	platform *domain.Platform,
) *domain.Publisher {
	return &domain.Publisher{
		Grading: publisher.NewGradingPublisher(cfg, platform.RabbitMq),
	}
}

func startConsumer(
	logger *zap.Logger,
	platform *domain.Platform,
	usecase *domain.Usecase,
) {
	if err := consumer.NewGradingConsumer(
		logger,
		platform.RabbitMq,
		platform.WebSocketHub,
		platform.InfluxDb,
		usecase.Assignment,
	); err != nil {
		logger.Fatal("Cannot start grading consumer", zap.Error(err))
	}
}
