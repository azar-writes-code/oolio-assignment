package products

import (
	"fmt"

	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	models_products "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/jackc/pgx/v5/pgtype"
)

func toDBProductParams(p *models_products.Product) (db.CreateProductParams, error) {
	var price pgtype.Numeric
	if err := price.Scan(fmt.Sprintf("%f", p.Price)); err != nil {
		return db.CreateProductParams{}, err
	}
	return db.CreateProductParams{
		Name:        p.Name,
		Description: p.Description,
		Price:       price,
		Category:    p.Category,
		Image:       p.Image,
		Stock:       p.Stock,
	}, nil
}

func toProduct(row db.CreateProductRow) db.Product {
	return db.Product{
		ID:          row.ID,
		Name:        row.Name,
		Description: row.Description,
		Price:       row.Price,
		Category:    row.Category,
		Image:       row.Image,
		Stock:       row.Stock,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}
