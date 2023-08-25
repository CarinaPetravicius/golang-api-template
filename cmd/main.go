package main

import (
	"github.com/go-playground/validator/v10"
	"golang-api-template/config"
	"golang-api-template/internal/adapters/api"
	"golang-api-template/internal/adapters/api/middleware"
)

func main() {
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	configs := config.LoadConfigFile(logger)

	database := config.NewDatabaseConnection(logger, configs.DB)
	defer config.CloseDatabaseConnection(database)

	prometheusMetrics := middleware.NewPrometheusMiddleware(configs.Service.Name)

	// Config Http Routers and Controllers
	router := api.NewHTTPRouter(prometheusMetrics)
	validate := validator.New()
	logger.Infof("Config validator: %v", validate)
	api.NewHealthCheckController(router, prometheusMetrics)

	config.StartHttpServer(logger, configs.Server, router)
}
