package api

import (
	"encoding/json"
	serverMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"golang-api-template/adapters/api/dto"
	"golang-api-template/adapters/api/middleware"
	"golang-api-template/core/domain"
	"net/http"
)

// ProductController controller for product API
type ProductController struct {
	log           *zap.SugaredLogger
	counterMetric prometheus.Counter
	validate      *validator.Validate
}

// NewProductController Create a new http product controller API
func NewProductController(httpRouter *HTTPRouter, log *zap.SugaredLogger, validator *validator.Validate, prometheusRegistry *middleware.CustomMetricRegistry) {
	controller := &ProductController{
		log:      log,
		validate: validator,
		counterMetric: promauto.With(prometheusRegistry).NewCounter(prometheus.CounterOpts{
			Name: "products_reqs_total",
			Help: "The total number of request for products endpoints",
		}),
	}

	httpRouter.Router.Post("/v1/product", controller.createProduct)
}

// createProduct create the product
func (p *ProductController) createProduct(writer http.ResponseWriter, request *http.Request) {
	p.counterMetric.Inc()
	traceID := request.Context().Value(serverMiddleware.RequestIDKey)

	productRequest := domain.Product{}
	err := json.NewDecoder(request.Body).Decode(&productRequest)
	if err != nil {
		p.log.With("traceId", traceID).Errorf("Error to parsing the product payload body. Maformed: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	_ = p.validate.RegisterValidation("not_blank", validators.NotBlank)
	err = p.validate.Struct(productRequest)
	if err != nil {
		p.log.With("traceId", traceID).Errorf("Product validation error: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	dto.RenderResponse(request.Context(), writer, http.StatusCreated, dto.DefaultResponse(http.StatusText(http.StatusCreated), ""))
}
