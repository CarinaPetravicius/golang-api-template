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

	productModel := domain.FromProductToProductModel(request, username)

	_, err = ps.productRepository.Create(ctx, productModel)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	}

	ps.log.With("traceId", traceID).Infof("The productID %s was created with success", productModel.ID)
	return &domain.ProductResponse{ID: productModel.ID}, nil
}

// GetProduct get the product by id
func (ps *ProductService) GetProduct(ctx context.Context, productID, traceID string) (*domain.ProductResponse, error) {
	productModel, err := ps.productRepository.GetProductById(ctx, productID)
	if err != nil {
		ps.log.With("traceId", traceID).Errorf("Internal server error to get the product: %v", err)
		return nil, custom_error.New(http.StatusInternalServerError, "internal server error")
	}
	if productModel == nil {
		ps.log.With("traceId", traceID).Errorf("Product not found")
		return nil, custom_error.New(http.StatusNotFound, "not found")
	}

	ps.log.With("traceId", traceID).Infof("The productID %s was found with success", productModel.ID)
	return domain.FromProductModelToProductResponse(productModel), nil
}
