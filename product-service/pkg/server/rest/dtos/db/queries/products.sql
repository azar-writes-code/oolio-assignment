-- name: CreateProduct :one
INSERT INTO products (name, description, price, image, category, stock)
VALUES (@name, @description, @price, @image, @category, @stock)
RETURNING id, name, description, price, image, category, stock, created_at, updated_at;

-- name: CreateManyProducts :many
INSERT INTO products (name, description, price, image, category, stock)
VALUES (
    unnest(@names::text[]),
    unnest(@descriptions::text[]),
    unnest(@prices::decimal[]),
    unnest(@images::product_image[]),
    string_to_array(unnest(@categories::text[]), '|||'),
    unnest(@stocks::int[])
)
RETURNING *;




-- name: GetProductByID :one
SELECT * FROM products
WHERE id = @id;

-- name: GetProductForUpdate :one
SELECT * FROM products
WHERE id = @id
FOR UPDATE;

-- name: ListProducts :many
SELECT * FROM products
ORDER BY id ASC
LIMIT @page_size OFFSET @page_offset;

-- name: UpdateProduct :one
UPDATE products SET name = COALESCE(@name, name), description = COALESCE(@description, description), price = COALESCE(@price, price), category = COALESCE(@category, category), image = COALESCE(@image, image), stock = COALESCE(@stock, stock), updated_at = NOW() WHERE id = @id RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = @id;

-- name: DecrementProductStock :execrows
UPDATE products
        SET stock = stock - @quantity
WHERE id = @id AND stock >= @quantity;

-- name: CountProducts :one
SELECT COUNT(*) FROM products;
