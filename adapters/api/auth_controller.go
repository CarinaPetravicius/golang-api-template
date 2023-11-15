package api

import (
	"encoding/json"
	serverMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"go.uber.org/zap"
	"golang-api-template/adapters/api/dto"
	"golang-api-template/adapters/api/router"
	"golang-api-template/core/domain"
	"golang-api-template/core/ports"
	"net/http"
)

// AuthController authentication controller
type AuthController struct {
	log      *zap.SugaredLogger
	validate *validator.Validate
	service  ports.IAuthService
}

// NewAuthController new authentication controller
func NewAuthController(httpRouter *router.HTTPRouter, log *zap.SugaredLogger, validator *validator.Validate, service ports.IAuthService) {
	controller := &AuthController{
		log:      log,
		validate: validator,
		service:  service,
	}

	httpRouter.Router.Post("/v1/sts/token", controller.createToken)
}

// createToken create a simple access token for tests
func (ac *AuthController) createToken(writer http.ResponseWriter, request *http.Request) {
	traceID := request.Context().Value(serverMiddleware.RequestIDKey).(string)
	auth := &domain.Auth{}

	err := json.NewDecoder(request.Body).Decode(auth)
	if err != nil {
		ac.log.With("traceId", traceID).Errorf("Error to parsing the authentication payload body. Maformed: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	_ = ac.validate.RegisterValidation("not_blank", validators.NotBlank)
	err = ac.validate.Struct(auth)
	if err != nil {
		ac.log.With("traceId", traceID).Errorf("Authentication validation error: %v", err)
		dto.RenderErrorResponse(request.Context(), writer, http.StatusBadRequest, err)
		return
	}

	tokenString, err := ac.service.CreateOauthToken(auth, traceID)
	if err != nil {
		dto.RenderErrorResponse(request.Context(), writer, http.StatusInternalServerError, err)
		return
	}

	dto.RenderResponse(request.Context(), writer, http.StatusOK, tokenString)
}
