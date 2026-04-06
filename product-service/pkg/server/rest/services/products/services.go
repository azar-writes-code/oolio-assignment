package products

import (
	"context"

	"github.com/azar-writes-code/oolio-products-backend/config"
	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
)

type ProductDAO interface {
	CountProducts(ctx context.Context) (int64, error)
	CreateManyProducts(ctx context.Context, arg db.CreateManyProductsParams) ([]db.Product, error)
	CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.CreateProductRow, error)
	DeleteProduct(ctx context.Context, id int32) error
	GetProductByID(ctx context.Context, id int32) (db.Product, error)
	ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Product, error)
	UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error)
}

type Service struct {
	cfg        *config.Config
	productDao ProductDAO
}

func NewService(cfg *config.Config, productDao ProductDAO) *Service {
	return &Service{
		cfg:        cfg,
		productDao: productDao,
	}
}
