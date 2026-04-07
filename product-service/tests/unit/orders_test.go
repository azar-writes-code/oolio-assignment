package unit

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/order"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/orders"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pashagolub/pgxmock/v3"
)

type mockCouponService struct {
	valid bool
	err   error
}

func (m *mockCouponService) Validate(code string) (bool, error) {
	return m.valid, m.err
}

func TestCreateOrder(t *testing.T) {
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		couponSvc := &mockCouponService{valid: true}
		svc := orders.NewService(mockPool, nil, couponSvc)

		mockPool.ExpectBegin()

		img := products.Image{
			Thumbnail: "http://example.com/t.jpg",
			Mobile:    "http://example.com/m.jpg",
			Tablet:    "http://example.com/ta.jpg",
			Desktop:   "http://example.com/d.jpg",
		}

		// GetProductForUpdate
		mockPool.ExpectQuery("SELECT .* FROM products WHERE id = \\$1 FOR UPDATE").
			WithArgs(int32(1)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "price", "category", "image", "stock", "created_at", "updated_at"}).
				AddRow(int32(1), "Product 1", nil, pgtype.Numeric{Int: big.NewInt(1000), Exp: -2, Valid: true}, []string{"cat1"}, img, int32(10), time.Now(), time.Now()))

		// CreateOrder
		mockPool.ExpectQuery("INSERT INTO orders").
			WithArgs(pgxmock.AnyArg(), pgxmock.AnyArg(), "completed").
			WillReturnRows(pgxmock.NewRows([]string{"id", "coupon_code", "total_amount", "status", "created_at", "updated_at"}).
				AddRow(pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, nil, pgtype.Numeric{Valid: true}, "completed", time.Now(), time.Now()))

		// DecrementProductStock
		mockPool.ExpectExec("UPDATE products").
			WithArgs(int32(1), int32(1)).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		// CreateOrderItem
		mockPool.ExpectExec("INSERT INTO order_items").
			WithArgs(pgtype.UUID{Bytes: [16]byte{1}, Valid: true}, int32(1), int32(1)).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		mockPool.ExpectCommit()

		req := order.OrderReq{
			Items: []order.OrderItemReq{
				{ProductID: 1, Quantity: 1},
			},
		}

		res, err := svc.CreateOrder(ctx, req)
		if err != nil {
			t.Fatalf("CreateOrder failed: %v", err)
		}
		if res.TotalPrice != 10.0 {
			t.Errorf("Expected total price 10.0, got %f", res.TotalPrice)
		}

		if err := mockPool.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("InsufficientStock", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		couponSvc := &mockCouponService{valid: true}
		svc := orders.NewService(mockPool, nil, couponSvc)

		mockPool.ExpectBegin()

		img := products.Image{}

		mockPool.ExpectQuery("SELECT .* FROM products WHERE id = \\$1 FOR UPDATE").
			WithArgs(int32(1)).
			WillReturnRows(pgxmock.NewRows([]string{"id", "name", "description", "price", "category", "image", "stock", "created_at", "updated_at"}).
				AddRow(int32(1), "Product 1", nil, pgtype.Numeric{Int: big.NewInt(1000), Exp: -2, Valid: true}, []string{"cat1"}, img, int32(0), time.Now(), time.Now()))

		mockPool.ExpectRollback()

		req := order.OrderReq{
			Items: []order.OrderItemReq{
				{ProductID: 1, Quantity: 1},
			},
		}

		_, err = svc.CreateOrder(ctx, req)
		if err == nil {
			t.Fatal("Expected error due to insufficient stock, got nil")
		}

		if err := mockPool.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %s", err)
		}
	})

	t.Run("InvalidCoupon", func(t *testing.T) {
		mockPool, err := pgxmock.NewPool()
		if err != nil {
			t.Fatal(err)
		}
		defer mockPool.Close()

		couponSvc := &mockCouponService{valid: false}
		svc := orders.NewService(mockPool, nil, couponSvc)

		couponCode := "INVALID"
		req := order.OrderReq{
			Items: []order.OrderItemReq{
				{ProductID: 1, Quantity: 1},
			},
			CouponCode: &couponCode,
		}

		_, err = svc.CreateOrder(ctx, req)
		if err == nil {
			t.Fatal("Expected error due to invalid coupon, got nil")
		}
	})
}
