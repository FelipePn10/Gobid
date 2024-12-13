-- name: CreatedProduct :one
INSERT INTO products (
    product_name, seller_id, description, baseprice, auction_end
) VALUES ($1, $2, $3, $4, $5)
RETURNING id;

-- name: UpdateProduct :exec
UPDATE products
SET
    product_name = COALESCE($3, product_name),
    description = COALESCE($4, description),
    baseprice = COALESCE($5, baseprice),
    auction_end = COALESCE($6, auction_end)
WHERE id = $1 AND seller_id = $2;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1 AND seller_id = $2;

-- name: GetProductById :one
SELECT * FROM products
WHERE id = $1;
