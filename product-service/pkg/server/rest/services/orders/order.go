package orders

import (
	"context"
	"fmt"

	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/coupon"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/order"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type Pool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct {
	pool      Pool
	queries   db.Querier
	couponSvc coupon.CouponService
}

func NewService(pool Pool, queries db.Querier, couponSvc coupon.CouponService) *Service {
	return &Service{
		pool:      pool,
		queries:   queries,
		couponSvc: couponSvc,
	}
}

func (s *Service) CreateOrder(ctx context.Context, req order.OrderReq) (*order.OrderResponse, error) {
	// 1. Validate coupon if exists
	var couponValid bool
	if req.CouponCode != nil && *req.CouponCode != "" {
		valid, err := s.couponSvc.Validate(*req.CouponCode)
		if err != nil {
			return nil, err // Let the apperror propagate
		}
		couponValid = valid
		if !valid {
			return nil, apperrors.NewBadRequest("invalid coupon code: %s", *req.CouponCode)
		}
	}

	// 2. Start Transaction
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, apperrors.NewInternal("failed to start transaction", err)
	}
	defer tx.Rollback(ctx)

	// Cast queries to use current transaction
	qtx := db.New(tx)

	// 3. Fetch products, lock them, and calculate total
	var total float64
	productsList := make([]products.Product, 0, len(req.Items))
	for _, item := range req.Items {
		// Use GetProductForUpdate to acquire a pessimistic lock on the row
		p, err := qtx.GetProductForUpdate(ctx, item.ProductID)
		if err != nil {
			return nil, apperrors.NewNotFound("product %d not found or currently locked", item.ProductID)
		}

		// Immediate stock check while holding the lock
		if p.Stock < int32(item.Quantity) {
			return nil, apperrors.NewConflict("insufficient stock for product %d", item.ProductID)
		}

		// Convert pgtype.Numeric price to float64
		priceVal, _ := p.Price.Float64Value()
		total += priceVal.Float64 * float64(item.Quantity)

		// Map for response products
		productsList = append(productsList, products.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       priceVal.Float64,
			Category:    p.Category,
			Image:       p.Image,
			Stock:       p.Stock,
		})
	}

	// 4. Calculate discount (0 for now as requested)
	discountedTotal := total
	if couponValid {
		// Future: add discount calculation logic here
		discountedTotal = total - 0
	}

	// 5. Create Numeric total_amount
	var numericTotal pgtype.Numeric
	_ = numericTotal.Scan(fmt.Sprintf("%.2f", total))

	// 6. Insert Order
	orderParams := db.CreateOrderParams{
		CouponCode:  req.CouponCode,
		TotalAmount: numericTotal,
		Status:      "completed",
	}

	dbOrder, err := qtx.CreateOrder(ctx, orderParams)
	if err != nil {
		return nil, apperrors.NewInternal("failed to create order", err)
	}

	// 7. Insert Items & Update Stock
	respItems := make([]order.OrderItem, len(req.Items))
	for i, item := range req.Items {
		// 7a. Decrement stock atomically (lock is already held)
		rowsAffected, err := qtx.DecrementProductStock(ctx, db.DecrementProductStockParams{
			ID:       item.ProductID,
			Quantity: int32(item.Quantity),
		})
		if err != nil {
			return nil, apperrors.NewInternal("failed to decrement stock", err)
		}
		if rowsAffected == 0 {
			return nil, apperrors.NewConflict("insufficient stock for product %d during update", item.ProductID)
		}

		// 7b. Update the product stock in our response list to reflect the change
		productsList[i].Stock -= int32(item.Quantity)

		// 7c. Insert order item
		err = qtx.CreateOrderItem(ctx, db.CreateOrderItemParams{
			OrderID:   dbOrder.ID,
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
		})
		if err != nil {
			return nil, apperrors.NewInternal("failed to create order item", err)
		}
		respItems[i] = order.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	// 8. Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, apperrors.NewInternal("failed to commit transaction", err)
	}

	// 9. Prepare Response
	return &order.OrderResponse{
		ID:              uuid.UUID(dbOrder.ID.Bytes).String(),
		Status:          dbOrder.Status,
		Items:           respItems,
		Products:        productsList,
		CouponValid:     couponValid,
		TotalPrice:      total,
		DiscountedPrice: discountedTotal,
	}, nil
}
