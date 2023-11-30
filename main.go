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
	"github.com/codern-org/codern/platform/server"
	"github.com/codern-org/codern/repository"
	"github.com/codern-org/codern/usecase"
	"github.com/jmoiron/sqlx"
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
	usecase := initUsecase(cfg, platform, repository)

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
	platform.InfluxDb.Close()
	platform.MySql.Close()
	platform.SeaweedFs.Close()
	platform.RabbitMq.Close()

	logger.Info("Server was successful shutdown")
}

func initPlatform(cfg *config.Config, logger *zap.Logger) *domain.Platform {
	start := time.Now()
	influxdb, err := platform.NewInfluxDb(
		cfg.Client.InfluxDb.Url,
		cfg.Client.InfluxDb.Token,
		cfg.Client.InfluxDb.Org,
		cfg.Client.InfluxDb.Bucket,
	)
	if err != nil {
		logger.Fatal("Cannot create a InfluxDB connection", zap.Error(err))
	}
	logger.Info("Connected to InfluxDB", zap.String("connection_time", time.Since(start).String()))

	start = time.Now()
	mysql, err := platform.NewMySql(cfg.Client.MySql.Uri)
	if err != nil {
		logger.Fatal("Cannot open MySQL database connection", zap.Error(err))
	}
	logger.Info("Connected to MySQL", zap.String("connection_time", time.Since(start).String()))

	start = time.Now()
	seaweedfs, err := platform.NewSeaweedFs(
		cfg.Client.SeaweedFs.MasterUrl,
		cfg.Client.SeaweedFs.FilerUrls,
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

	return &domain.Platform{
		InfluxDb:  influxdb,
		MySql:     mysql,
		SeaweedFs: seaweedfs,
		RabbitMq:  rabbitmq,
	}
}

func initRepository(mysql *sqlx.DB) *domain.Repository {
	return &domain.Repository{
		Session:   repository.NewSessionRepository(mysql),
		User:      repository.NewUserRepository(mysql),
		Workspace: repository.NewWorkspaceRepository(mysql),
	}
}

func initUsecase(
	cfg *config.Config,
	platform *domain.Platform,
	repository *domain.Repository,
) *domain.Usecase {
	googleUsecase := usecase.NewGoogleUsecase(cfg)
	sessionUsecase := usecase.NewSessionUsecase(cfg, repository.Session)
	userUsecase := usecase.NewUserUsecase(repository.User)
	authUsecase := usecase.NewAuthUsecase(googleUsecase, sessionUsecase, userUsecase)
	workspaceUsecase := usecase.NewWorkspaceUsecase(
		cfg, platform.SeaweedFs, platform.RabbitMq, repository.Workspace,
	)

	return &domain.Usecase{
		Google:    googleUsecase,
		Session:   sessionUsecase,
		User:      userUsecase,
		Auth:      authUsecase,
		Workspace: workspaceUsecase,
	}
}
