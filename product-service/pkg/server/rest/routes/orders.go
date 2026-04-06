package routes

import (
	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/controllers/order"
	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	badgerRepo "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/badger"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/middleware"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/coupons"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/orders"
	"github.com/dgraph-io/badger/v4"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// OrderRoutes handles order placement and management.
func OrderRoutes(rg *gin.RouterGroup, cfg *config.Config, pool *pgxpool.Pool, badgerDB *badger.DB) {
	// 1. Initialize DB Queries
	queries := db.New(pool)

	// Initialize coupon service (Badger)
	couponRepo := badgerRepo.New(badgerDB)
	couponSvc := coupons.NewService(couponRepo)

	// 3. Initialize Order Service & Controller
	orderSvc := orders.NewService(pool, queries, couponSvc)
	orderCtrl := order.NewOrderController(orderSvc, couponSvc)

	// 4. Register Routes
	orders := rg.Group("/order")
	orders.Use(middleware.AuthMiddleware(cfg.Auth.ApiKey))
	{
		orders.POST("/", orderCtrl.CreateOrder)
		orders.POST("/validate-coupon", orderCtrl.ValidateCoupon)
	}
}
