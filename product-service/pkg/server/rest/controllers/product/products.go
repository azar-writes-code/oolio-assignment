package product

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models"
	models_products "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/gin-gonic/gin"
)

func (p *ProductController) CreateProduct(c *gin.Context) {
	var req models_products.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("failed to bind product data: %v", err)))
		return
	}

	if err := req.Validate(); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("validation failed: %v", err)))
		return
	}

	res, err := p.productSvc.CreateProduct(c.Request.Context(), &req)
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p *ProductController) GetProducts(c *gin.Context) {
	var pageParams models.PaginationParams
	if err := c.ShouldBindQuery(&pageParams); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewBadRequest("invalid pagination parameters: %v", err)))
		return
	}

	res, err := p.productSvc.ListProducts(c.Request.Context(), pageParams.Page, pageParams.PageSize)
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p *ProductController) CreateManyProducts(c *gin.Context) {
	var req []models_products.Product
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("failed to bind products data: %v", err)))
		return
	}

	// Validate all products first
	for _, product := range req {
		if err := product.Validate(); err != nil {
			c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("validation failed for product %s: %v", product.Name, err)))
			return
		}
	}

	// Logic for mapping moved to service
	res, err := p.productSvc.CreateMany(c.Request.Context(), req)
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (p *ProductController) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		slog.Error("Invalid product ID", "id", idStr, "error", err.Error())
		c.JSON(apperrors.MapToHTTP(apperrors.NewBadRequest("invalid product ID: %s", idStr)))
		return
	}

	res, err := p.productSvc.GetProductByID(c.Request.Context(), int32(id))
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	c.JSON(http.StatusOK, res)
}