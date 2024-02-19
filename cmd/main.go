package main

import (
	"context"
	"github.com/go-playground/validator/v10"
	"golang-api-template/adapters/api"
	middleware2 "golang-api-template/adapters/api/middleware"
	"golang-api-template/adapters/api/router"
	"golang-api-template/adapters/kafka"
	"golang-api-template/adapters/opa"
	"golang-api-template/adapters/repository/products"
	"golang-api-template/config"
	"golang-api-template/core/services"
)

func main() {
	logger := config.NewLogger()
	defer config.CloseLogger(logger)

	configs := config.LoadConfigFile(logger)

	database := config.NewDatabaseConnection(logger, configs.DB)
	defer config.CloseDatabaseConnection(database)

	// Opa Policies
	policies := opa.NewPolicyService(configs.Policies.Path, logger)

	// Repositories
	productsRepository := products.NewProductRepository(database)

	// Config Domain Services
	productService := services.NewProductService(logger, productsRepository)
	authService := services.NewAuthService(logger, configs.Oauth)

	prometheusMetrics := middleware2.NewPrometheusMiddleware(configs.Service.Name)
	jwtHandler := middleware2.NewJWTHandler(logger, authService)

	// Config Http Routers and Controllers
	route := router.NewHTTPRouter(prometheusMetrics)
	valid := validator.New()
	api.NewHealthCheckController(route, prometheusMetrics)
	api.NewAuthController(route, logger, valid, authService)
	api.NewProductController(route, logger, valid, prometheusMetrics, productService, jwtHandler, policies)

	// Start Kafka with a new context
	ctx := context.Background()
	producer := kafka.NewKafkaProducer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Producer))
	defer producer.Close()
	kafka.CreateKafkaTopics(logger, configs.Kafka, ctx, config.NewKafkaConfigMap(logger, configs.Kafka, config.Topic))
	consumer := kafka.NewKafkaConsumer(logger, config.NewKafkaConfigMap(logger, configs.Kafka, config.Consumer))
	defer kafka.CloseConsumer(consumer)

	config.StartHttpServer(logger, configs.Server, route)
}
