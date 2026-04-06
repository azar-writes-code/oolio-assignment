package order

import "github.com/azar-writes-code/oolio-products-backend/pkg/server/rest/models/products"

type OrderItemReq struct {
	ProductID int32 `json:"productId" binding:"required"`
	Quantity  int   `json:"quantity" binding:"required,gt=0"`
}

type OrderReq struct {
	CouponCode *string        `json:"couponCode" binding:"omitempty"`
	Items      []OrderItemReq `json:"items" binding:"required,dive"`
}

type OrderItem struct {
	ProductID int32 `json:"productId"`
	Quantity  int   `json:"quantity"`
}

type OrderResponse struct {
	ID              string             `json:"id"`
	Status          string             `json:"status"`
	Items           []OrderItem        `json:"items"`
	Products        []products.Product `json:"products"`
	CouponValid     bool               `json:"couponValid"`
	TotalPrice      float64            `json:"totalPrice"`
	DiscountedPrice float64            `json:"discountedPrice"`
}


type ValidateCouponReq struct {
	CouponCode string `json:"couponCode" binding:"required"`
}

type ValidateCouponResponse struct {
	IsValid bool `json:"isValid"`
}