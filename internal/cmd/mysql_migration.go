package main

import (
	"flag"

	"github.com/codern-org/codern/internal/config"
	"github.com/codern-org/codern/internal/logger"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger := logger.NewLogger()

	// Load configuration file
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "path to a config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		logger.Fatal("Cannot load a config file", zap.Error(err))
	}
	logger.Info("Configuration file loaded successfully")

	// Initialize migrator
	m, err := migrate.New(
		"file:///workspace/other/db/migrations/",
		"mysql://"+cfg.Client.MySql.Uri,
	)
	if err != nil {
		logger.Error("Cannot run mysql migration", zap.Error(err))
	}

	err = m.Up()
	if err != nil {
		logger.Error("Cannot run mysql migration", zap.Error(err))
	}

	logger.Info("MySQL migration done")
}
