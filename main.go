package main

import (
	"flag"
	"time"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/constant"
	"github.com/codern-org/codern/internal/logger"
	"github.com/codern-org/codern/platform"
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

	// Load configuration file
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "path to a config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal("Cannot load a config file", zap.Error(err))
	}
	logger.Info("Configuration file loaded successfully")

	// Initialize databases
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

	// Initialize HTTP server
	fiber := platform.NewFiberServer(cfg, logger, influxdb, mysql)
	fiber.Start()
}
