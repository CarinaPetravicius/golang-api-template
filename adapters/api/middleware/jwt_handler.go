package middleware

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang-api-template/adapters/api/dto"
	"golang-api-template/core/domain"
	"net/http"
	"strings"
	"time"
)

// JWTVerify jwt verify token
type JWTVerify struct {
	log *zap.SugaredLogger
}

// NewJWTHandler create new JWT verify handler
func NewJWTHandler(log *zap.SugaredLogger) *JWTVerify {
	return &JWTVerify{
		log: log,
	}
}

// JWTVerifyHandler handler jwt verify
func (jw *JWTVerify) JWTVerifyHandler() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, err := jw.parseTokenFromRequest(r)
			if err != nil {
				jw.log.Errorf("JWT parsing failed: %v", err)
				dto.RenderErrorResponse(r.Context(), w, http.StatusUnauthorized, err)
				return
			}

			context.WithValue(r.Context(), domain.Claims, claims)
			next.ServeHTTP(w, r)
		})
	}
}

func (jw *JWTVerify) parseTokenFromRequest(r *http.Request) (*domain.AuthClaims, error) {
	header := r.Header.Get("Authorization")
	if len(header) == 0 {
		jw.log.Error("no security header")
		return nil, errors.New("no security header")
	}

	tokenString := strings.Split(header, "Bearer ")
	if len(tokenString) == 0 {
		jw.log.Error("no security header token")
		return nil, errors.New("no security header token")
	}

	token, err := jwt.Parse(tokenString[2], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return domain.SecretKey, nil
	})

	if err != nil || !token.Valid {
		jw.log.Errorf("invalid token: %v", err)
		return nil, errors.New("invalid token")
	}

	issuer, err := token.Claims.GetIssuer()
	if err != nil || issuer != domain.Issuer {
		jw.log.Errorf("invalid token issuer: %v", err)
		return nil, errors.New("invalid token issuer")
	}

	expiration, err := token.Claims.GetExpirationTime()
	if err != nil || expiration.Time.Before(time.Now()) {
		jw.log.Errorf("token expired: %v", err)
		return nil, errors.New("token expired")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return &domain.AuthClaims{
			Username: claims["username"].(string),
			Role:     claims["role"].(string),
		}, nil
	} else {
		jw.log.Errorf("claims not found")
		return nil, errors.New("claims not found")
	}
}
