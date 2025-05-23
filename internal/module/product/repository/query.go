package repository

const (
	queryInsertNewProduct = `
		INSERT INTO products 
		(
			name, 
			description, 
			image_url,
			stock, 
			is_active,
			product_category_id,
			real_price,
			recipe_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		RETURNING id, name, description, image_url, stock, is_active, product_category_id, real_price, recipe_id, created_at
	`

	queryUpdateProduct = `
		UPDATE products
		SET
			name = ?,
			description = ?,
			image_url = ?,
			stock = ?,
			is_active = ?,
			product_category_id = ?,
			real_price = ?,
			recipe_id = ?
		WHERE id = ? AND deleted_at IS NULL
		RETURNING id, name, description, image_url, stock, is_active, product_category_id, real_price, recipe_id
	`

	queryCountProductByID = `
		SELECT COUNT(*)
		FROM products
		WHERE id = ? AND deleted_at IS NULL
	`

	queryFindProductByID = `
		SELECT
			p.id,
			p.name,
			p.description,
			p.image_url,
			p.stock,
			p.is_active,
			p.product_category_id,
			p.real_price,
			pc.name AS product_category_name
		FROM products p
		JOIN product_categories pc ON p.product_category_id = pc.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`

	queryFindAllProduct = `
		SELECT
			p.id,
			p.name,
			p.description,
			p.image_url,
			p.stock,
			p.is_active,
			p.product_category_id,
			p.real_price,
			pc.name AS product_category_name
		FROM products p
		JOIN product_categories pc ON p.product_category_id = pc.id
		WHERE p.deleted_at IS NULL
		AND pc.deleted_at IS NULL
		AND (
			p.name ILIKE '%' || ? || '%'
			OR p.description ILIKE '%' || ? || '%'
		)
		ORDER BY p.created_at DESC, p.id DESC
		LIMIT ? OFFSET ?
	`

	queryCountFindAllProduct = `
		SELECT COUNT(1)
		FROM products p
		JOIN product_categories pc ON p.product_category_id = pc.id
		WHERE
			p.deleted_at IS NULL
			AND (
				p.name ILIKE '%' || ? || '%'
				OR p.description ILIKE '%' || ? || '%'
			)
	`

	qyerySoftDeleteArticleByID = `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`
)
