package main

import (
	"github.com/go-playground/validator/v10"
	"golang-api-template/config"
)

func main() {
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	configs := config.LoadConfigFile(logger)

	database := config.NewDatabaseConnection(logger, configs.DB)
	defer config.CloseDatabaseConnection(database)

	validate := validator.New()
	logger.Infof("Config validator: %v", validate)

	config.StartHttpServer(logger, configs.Server)
}
