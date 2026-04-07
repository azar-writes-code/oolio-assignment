package unit

import (
	"context"
	"fmt"
	"testing"

	"github.com/azar-writes-code/oolio-products-backend/config"
	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	models_products "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/products"
)

type mockProductDAO struct {
	products map[int32]db.Product
	nextID   int32
}

func (m *mockProductDAO) CountProducts(ctx context.Context) (int64, error) {
	return int64(len(m.products)), nil
}

func (m *mockProductDAO) CreateProduct(ctx context.Context, arg db.CreateProductParams) (db.CreateProductRow, error) {
	m.nextID++
	id := m.nextID
	p := db.Product{
		ID:          id,
		Name:        arg.Name,
		Description: arg.Description,
		Price:       arg.Price,
		Category:    arg.Category,
		Image:       arg.Image,
		Stock:       arg.Stock,
	}
	m.products[id] = p
	return db.CreateProductRow{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Category:    p.Category,
		Image:       p.Image,
		Stock:       p.Stock,
	}, nil
}

func (m *mockProductDAO) CreateManyProducts(ctx context.Context, arg db.CreateManyProductsParams) ([]db.Product, error) {
	var result []db.Product
	for i := range arg.Names {
		m.nextID++
		id := m.nextID
		p := db.Product{
			ID:          id,
			Name:        arg.Names[i],
			Description: &arg.Descriptions[i],
			Price:       arg.Prices[i],
			Category:    []string{arg.Categories[i]},
			Stock:       arg.Stocks[i],
		}
		m.products[id] = p
		result = append(result, p)
	}
	return result, nil
}

func (m *mockProductDAO) GetProductByID(ctx context.Context, id int32) (db.Product, error) {
	p, ok := m.products[id]
	if !ok {
		return db.Product{}, fmt.Errorf("not found")
	}
	return p, nil
}

func (m *mockProductDAO) DeleteProduct(ctx context.Context, id int32) error {
	if _, ok := m.products[id]; !ok {
		return fmt.Errorf("not found")
	}
	delete(m.products, id)
	return nil
}

func (m *mockProductDAO) ListProducts(ctx context.Context, arg db.ListProductsParams) ([]db.Product, error) {
	var res []db.Product
	for _, p := range m.products {
		res = append(res, p)
	}
	return res, nil
}

func (m *mockProductDAO) UpdateProduct(ctx context.Context, arg db.UpdateProductParams) (db.Product, error) {
	p, ok := m.products[arg.ID]
	if !ok {
		return db.Product{}, fmt.Errorf("not found")
	}
	p.Name = arg.Name
	p.Price = arg.Price
	m.products[arg.ID] = p
	return p, nil
}

func TestProductService(t *testing.T) {
	dao := &mockProductDAO{products: make(map[int32]db.Product)}
	svc := products.NewService(&config.Config{}, dao)
	ctx := context.Background()

	t.Run("CreateProduct", func(t *testing.T) {
		desc := "test description"
		req := &models_products.Product{
			Name:        "Test Product",
			Description: &desc,
			Price:       99.99,
			Category:    []string{"electronics"},
			Stock:       10,
		}
		res, err := svc.CreateProduct(ctx, req)
		if err != nil {
			t.Fatalf("CreateProduct failed: %v", err)
		}
		if res.Name != req.Name {
			t.Errorf("Expected name %s, got %s", req.Name, res.Name)
		}
	})

	t.Run("GetProductByID", func(t *testing.T) {
		res, err := svc.GetProductByID(ctx, 1)
		if err != nil {
			t.Fatalf("GetProductByID failed: %v", err)
		}
		if res.ID != 1 {
			t.Errorf("Expected ID 1, got %d", res.ID)
		}
	})

	t.Run("ListProducts", func(t *testing.T) {
		res, err := svc.ListProducts(ctx, 1, 10)
		if err != nil {
			t.Fatalf("ListProducts failed: %v", err)
		}
		if res.TotalCount != 1 {
			t.Errorf("Expected 1 product, got %d", res.TotalCount)
		}
	})

	t.Run("DeleteProduct", func(t *testing.T) {
		err := svc.DeleteProduct(ctx, 1)
		if err != nil {
			t.Fatalf("DeleteProduct failed: %v", err)
		}
		_, err = svc.GetProductByID(ctx, 1)
		if err == nil {
			t.Errorf("Expected product 1 to be deleted")
		}
	})
}
