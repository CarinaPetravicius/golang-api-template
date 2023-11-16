package api

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	serverMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"golang-api-template/adapters/api/dto"
	"golang-api-template/adapters/api/middleware"
	"golang-api-template/adapters/api/router"
	"golang-api-template/core/domain"
	"golang-api-template/core/ports"
	"net/http"
)

// ProductController controller for product API
type ProductController struct {
	log           *zap.SugaredLogger
	counterMetric prometheus.Counter
	validate      *validator.Validate
	service       ports.IProductService
	jwtVerify     *middleware.JWTVerify
}

// NewProductController Create a new http product controller API
func NewProductController(httpRouter *router.HTTPRouter, log *zap.SugaredLogger, validator *validator.Validate, prometheusRegistry *middleware.CustomMetricRegistry,
	service ports.IProductService, jwtVerify *middleware.JWTVerify) {
	controller := &ProductController{
		log:       log,
		validate:  validator,
		service:   service,
		jwtVerify: jwtVerify,
		counterMetric: promauto.With(prometheusRegistry).NewCounter(prometheus.CounterOpts{
			Name: "products_reqs_total",
			Help: "The total number of request for products endpoints",
		}),
	}

	httpRouter.Router.Group(func(r chi.Router) {
		r.Use(controller.jwtVerify.JWTVerifyHandler())
		r.Post("/v1/product", controller.createProduct)
	})
}

// createProduct create the product
func (pc *ProductController) createProduct(writer http.ResponseWriter, request *http.Request) {
	pc.counterMetric.Inc()
	traceID := request.Context().Value(serverMiddleware.RequestIDKey).(string)
	claims := request.Context().Value(domain.ClaimsKey).(domain.AuthClaims)
	pc.log.With("traceId", traceID).Infof("User %v is creating a product.", claims.Username)

	productRequest := &domain.Product{}
	err := json.NewDecoder(request.Body).Decode(productRequest)
	if err != nil {
		pc.log.With("traceId", traceID).Errorf("Error to parsing the product payload body. Maformed: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	_ = pc.validate.RegisterValidation("not_blank", validators.NotBlank)
	err = pc.validate.Struct(productRequest)
	if err != nil {
		pc.log.With("traceId", traceID).Errorf("Product validation error: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	response := pc.service.CreateProduct(request.Context(), productRequest, traceID)
	dto.RenderResponse(request.Context(), writer, http.StatusCreated, response)
}
