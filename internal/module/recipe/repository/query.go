package repository

const (
	queryInsertNewRecipe = `
		INSERT INTO recipes (capital_price) VALUES (?) RETURNING id, capital_price
	`

	queryCountRecipe = `
		SELECT COUNT(*) FROM recipes
	`

	queryFindListRecipeByProductIDs = `
		SELECT id, product_id, capital_price
		FROM recipes
		WHERE product_id = ANY($1)
	`

	queryFindRecipesByIDs = `
		SELECT id, capital_price
		FROM recipes
		WHERE id = ANY($1::int[])
		AND deleted_at IS NULL
	`

	queryUpdateRecipeCapitalPrice = `
		UPDATE recipes
        SET 
			capital_price = $2, 
			updated_at = NOW()
        WHERE id = $1
	`

	queryFindDetailRecipes = `
		SELECT
			r.id            AS id,
			r.capital_price AS capital_price,
			p.id            AS product_id,
			p.name          AS product_name,
			p.stock         AS product_stock
		FROM recipes r
		JOIN products p ON p.recipe_id = r.id
		WHERE r.id = $1 AND r.deleted_at IS NULL
	`

	queryFindDetailRecipesWIthIngredients = `
		SELECT
			i.id              AS id,
			i.unit            AS unit,
			i.cost            AS cost,
			CAST(i.required_stock AS INTEGER)    AS required_stock,
			i.recipe_id       AS recipe_id,
			s.id              AS stock_id,
			s.name            AS stock_name
		FROM ingredients i
		JOIN ingredient_stocks s ON s.id = i.ingredient_stock_id
		WHERE i.recipe_id = $1
	`

	queryFindIngredientsWithStock = `
		SELECT ir.recipe_id, ir.ingredient_id, ir.required_stock,
               s.id AS stock_id, s.required_stock AS stock_amount
          FROM ingredient_recipes ir
          JOIN ingredient_stocks s ON s.ingredient_id = ir.ingredient_id
         WHERE ir.recipe_id IN (?)
	`
)
