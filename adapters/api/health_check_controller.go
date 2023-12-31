package api

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang-api-template/adapters/api/dto"
	"golang-api-template/adapters/api/middleware"
	"golang-api-template/adapters/api/router"
	"net/http"
)

// NewHealthCheckController create a new health check http controller
func NewHealthCheckController(httpRouter *router.HTTPRouter, prometheusRegistry *middleware.CustomMetricRegistry) {
	// health check endpoints for kubernetes
	httpRouter.Router.Get("/health/live", handleLivelinessCheck)
	httpRouter.Router.Get("/health/ready", handleReadinessCheck)
	// prometheus metrics endpoint
	httpRouter.Router.Get("/metrics", promhttp.HandlerFor(prometheusRegistry, promhttp.HandlerOpts{}).ServeHTTP)
}

func handleLivelinessCheck(writer http.ResponseWriter, reader *http.Request) {
	dto.RenderResponse(reader.Context(), writer, http.StatusOK, dto.DefaultResponse(http.StatusText(http.StatusOK), ""))
}

func handleReadinessCheck(writer http.ResponseWriter, reader *http.Request) {
	dto.RenderResponse(reader.Context(), writer, http.StatusOK, dto.DefaultResponse(http.StatusText(http.StatusOK), ""))
}
