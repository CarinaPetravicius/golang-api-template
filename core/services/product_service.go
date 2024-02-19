package services

import (
	"context"
	"go.uber.org/zap"
	"golang-api-template/adapters/custom_error"
	"golang-api-template/adapters/repository/products"
	"golang-api-template/core/domain"
	"net/http"
)

// ProductService product service
type ProductService struct {
	log               *zap.SugaredLogger
	productRepository products.IRepository
}

// NewProductService create new product service
func NewProductService(log *zap.SugaredLogger, productRepository products.IRepository) *ProductService {
	return &ProductService{
		log:               log,
		productRepository: productRepository,
	}
}

// CreateProduct service to create the product
func (ps *ProductService) CreateProduct(ctx context.Context, request *domain.Product, username, traceID string) (*domain.ProductResponse, error) {
	exist, err := ps.productRepository.ProductAlreadyExist(ctx,
		request.Name, request.UnitType, request.Unit, request.Brand, request.Color, request.Style)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	} else if exist {
		ps.log.With("traceId", traceID).Errorf("Product already exist")
		return nil, custom_error.New(http.StatusConflict, "already exist")
	}

	productDomain := domain.FromProductToProductModel(request, username)

	_, err = ps.productRepository.Create(ctx, productDomain)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	}

	return &domain.ProductResponse{ID: productDomain.ID}, nil
}
