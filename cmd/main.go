package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"golang-api-template/adapters/api"
	"golang-api-template/adapters/api/middleware"
	"golang-api-template/adapters/kafka"
	"golang-api-template/config"
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
	_ = validator.New()
	api.NewHealthCheckController(router, prometheusMetrics)

	// Start Kafka with a new context
	ctx := context.Background()
	kafkaConfigMap := config.NewKafkaConfigMap(logger, configs.Kafka)
	producer := kafka.NewKafkaProducer(logger, kafkaConfigMap)
	defer producer.Close()
	kafka.CreateKafkaTopics(logger, kafkaConfigMap, configs.Kafka, ctx)

	config.StartHttpServer(logger, configs.Server, router)
}
