package ports

import (
	"golang-api-template/core/domain"
)

// IAuthService auth service interface
type IAuthService interface {
	CreateOauthToken(request *domain.Auth, traceID string) (string, error)
	ParseOauthToken(token string) (*domain.AuthClaims, error)
}
