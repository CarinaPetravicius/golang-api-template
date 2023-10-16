package services

import (
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang-api-template/core/domain"
)

// ProductService product service
type ProductService struct {
	log *zap.SugaredLogger
}

// NewProductService create new product service
func NewProductService(log *zap.SugaredLogger) *ProductService {
	return &ProductService{
		log: log,
	}
}

// CreateProduct service to create the product
func (ps *ProductService) CreateProduct(ctx context.Context, request domain.Product, traceID string) *domain.ProductResponse {
	return &domain.ProductResponse{Id: uuid.NewString()}
}
