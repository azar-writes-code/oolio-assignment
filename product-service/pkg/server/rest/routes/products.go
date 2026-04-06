package routes

import (
	"github.com/azar-writes-code/oolio-products-backend/config"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/controllers/product"
	db "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/dtos/db/sqlc"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/middleware"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/products"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ProductsRoutes(router *gin.RouterGroup, cfg *config.Config, pool *pgxpool.Pool) {
	dbQuerier := db.New(pool)
	productSvc := products.NewService(cfg, dbQuerier)
	productController := product.NewProductController(cfg, productSvc)

	productsRoutes := router.Group("/products")
	productsRoutes.Use(middleware.AuthMiddleware(cfg.Auth.ApiKey))
	productsRoutes.POST("/create-product", productController.CreateProduct)
	productsRoutes.POST("/create-many-products", productController.CreateManyProducts)
	productsRoutes.GET("/", productController.GetProducts)
	productsRoutes.GET("/:id", productController.GetProductByID)
}
