package products

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models"
	models_products "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/jackc/pgx/v5/pgtype"
)

func (s *Service) CreateProduct(ctx context.Context, prod *models_products.Product) (db.Product, error) {
	arg, err := toDBProductParams(prod)
	if err != nil {
		slog.Error("Failed to convert product to DB format", "error", err.Error())
		return db.Product{}, err
	}
	row, err := s.productDao.CreateProduct(ctx, arg)
	if err != nil {
		return db.Product{}, apperrors.NewInternal("failed to create product", err)
	}

	return toProduct(row), nil
}

func (s *Service) CreateMany(ctx context.Context, req []models_products.Product) ([]db.Product, error) {
	var args db.CreateManyProductsParams
	for _, product := range req {
		var descStr string
		if product.Description != nil {
			descStr = *product.Description
		}

		var price pgtype.Numeric
		if err := price.Scan(fmt.Sprintf("%f", product.Price)); err != nil {
			return nil, apperrors.NewUnprocessable("invalid price for product %s: %v", product.Name, err)
		}

		args.Names = append(args.Names, product.Name)
		args.Descriptions = append(args.Descriptions, descStr)
		args.Prices = append(args.Prices, price)
		args.Stocks = append(args.Stocks, product.Stock)
		args.Images = append(args.Images, product.Image)
		args.Categories = append(args.Categories, strings.Join(product.Category, "|||"))
	}

	res, err := s.productDao.CreateManyProducts(ctx, args)
	if err != nil {
		return nil, apperrors.NewInternal("failed to create products", err)
	}
	return res, nil
}

func (s *Service) GetProductByID(ctx context.Context, id int32) (db.Product, error) {
	product, err := s.productDao.GetProductByID(ctx, id)
	if err != nil {
		// Assuming generic error means not found if it's the DAO
		return db.Product{}, apperrors.NewNotFound("product %d not found", id)
	}
	return product, nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int32) error {
	err := s.productDao.DeleteProduct(ctx, id)
	if err != nil {
		return apperrors.NewInternal(fmt.Sprintf("failed to delete product %d", id), err)
	}
	return nil
}

func (s *Service) ListProducts(ctx context.Context, page, pageSize int32) (*models.PaginatedResponse, error) {
	// Default values
	if page == 0 {
		page = 1
	}
	if pageSize == 0 {
		pageSize = 10
	}

	args := db.ListProductsParams{
		PageSize:   pageSize,
		PageOffset: (page - 1) * pageSize,
	}

	products, err := s.productDao.ListProducts(ctx, args)
	if err != nil {
		return nil, apperrors.NewInternal("failed to list products", err)
	}

	totalCount, err := s.productDao.CountProducts(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("failed to count products", err)
	}

	totalPages := (int32(totalCount) + pageSize - 1) / pageSize

	return &models.PaginatedResponse{
		Data:       products,
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
