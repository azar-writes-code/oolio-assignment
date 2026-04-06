-- name: CreateOrder :one
INSERT INTO orders (coupon_code, total_amount, status)
VALUES (@coupon_code, @total_amount, @status)
RETURNING id, coupon_code, total_amount, status, created_at, updated_at;

-- name: CreateOrderItem :exec
INSERT INTO order_items (order_id, product_id, quantity)
VALUES (@order_id, @product_id, @quantity);

-- name: GetOrderByID :one  
SELECT * FROM orders WHERE id = @id;

-- name: ListOrderItems :many
SELECT product_id, quantity FROM order_items WHERE order_id = @order_id;

-- name: ListOrders :many
SELECT * FROM orders
ORDER BY created_at DESC
LIMIT @page_size OFFSET @page_offset;
