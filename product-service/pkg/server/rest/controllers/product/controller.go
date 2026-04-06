package product

import (
	"context"

	"github.com/azar-writes-code/oolio-products-backend/config"
	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models"
	models_products "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
)

type ProductService interface {
	CreateProduct(ctx context.Context, arg *models_products.Product) (db.Product, error)
	CreateMany(ctx context.Context, arg []models_products.Product) ([]db.Product, error)
	GetProductByID(ctx context.Context, id int32) (db.Product, error)
	DeleteProduct(ctx context.Context, id int32) error
	ListProducts(ctx context.Context, page, pageSize int32) (*models.PaginatedResponse, error)
}

type ProductController struct {
	config     *config.Config
	productSvc ProductService
}

func NewProductController(config *config.Config, productSvc ProductService) *ProductController {
	return &ProductController{
		config:     config,
		productSvc: productSvc,
	}
}
