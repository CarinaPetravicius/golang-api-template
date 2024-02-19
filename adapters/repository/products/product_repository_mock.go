package products

import (
	"context"
	"golang-api-template/core/domain"
)

// ProductRepositoryMock product repository mock
type ProductRepositoryMock struct{}

var (
	createFunc          func(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error)
	productAlreadyExist func(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error)
)

// Create is the repository mock for Create func
func (pr *ProductRepositoryMock) Create(ctx context.Context, model *domain.ProductModel) (*domain.ProductModel, error) {
	return createFunc(ctx, model)
}

// ProductAlreadyExist is the repository mock for ProductAlreadyExist func
func (pr *ProductRepositoryMock) ProductAlreadyExist(ctx context.Context, name, unitType, unit, brand, color, style string) (bool, error) {
	return productAlreadyExist(ctx, name, unitType, unit, brand, color, style)
}
