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
	valid := validator.New()
	api.NewHealthCheckController(router, prometheusMetrics)
	api.NewProductController(router, logger, valid, prometheusMetrics)

	// Start Kafka with a new context
	ctx := context.Background()
	producer := kafka.NewKafkaProducer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Producer))
	defer producer.Close()
	kafka.CreateKafkaTopics(logger, configs.Kafka, ctx, config.NewKafkaConfigMap(logger, configs.Kafka, config.Topic))
	consumer := kafka.NewKafkaConsumer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Consumer))
	defer kafka.CloseConsumer(consumer)

	config.StartHttpServer(logger, configs.Server, router)
}
