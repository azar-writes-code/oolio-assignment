package order

import (
	"net/http"

	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/coupon"
	orderModel "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/order"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/services/orders"
	"github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/utils/apperrors"
	"github.com/gin-gonic/gin"
)

type OrderController struct {
	svc       *orders.Service
	couponSvc coupon.CouponService
}

func NewOrderController(svc *orders.Service, couponSvc coupon.CouponService) *OrderController {
	return &OrderController{svc: svc, couponSvc: couponSvc}
}

func (ctrl *OrderController) CreateOrder(c *gin.Context) {
	// 1. Bind JSON
	var req orderModel.OrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("Validation exception: %v", err)))
		return
	}

	// 2. Call Service
	resp, err := ctrl.svc.CreateOrder(c.Request.Context(), req)
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	// 3. Return success
	c.JSON(http.StatusOK, resp)
}

func (ctrl *OrderController) ValidateCoupon(c *gin.Context) {
	// 1. Bind JSON
	var req orderModel.ValidateCouponReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(apperrors.MapToHTTP(apperrors.NewUnprocessable("Validation exception: %v", err)))
		return
	}

	// 2. Call Service
	resp, err := ctrl.couponSvc.Validate(req.CouponCode)
	if err != nil {
		c.JSON(apperrors.MapToHTTP(err))
		return
	}

	// 3. Return success
	c.JSON(http.StatusOK, gin.H{
		"isValid":    resp,
		"couponCode": req.CouponCode,
	})
}
