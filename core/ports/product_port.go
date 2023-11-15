package ports

import (
	"context"
	"golang-api-template/core/domain"
)

// IProductService product service interface
type IProductService interface {
	CreateProduct(ctx context.Context, request *domain.Product, traceID string) *domain.ProductResponse
}
