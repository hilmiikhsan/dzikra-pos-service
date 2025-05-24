package repository

const (
	queryFindIngredientByRecipeIDs = `
		SELECT id, unit, cost, required_stock, recipe_id, ingredient_stock_id
		FROM ingredients
		WHERE recipe_id = ANY($1::int[]) AND deleted_at IS NULL
	`

	queryDeleteIngredientByRecipeID = `
		DELETE FROM ingredients
        WHERE recipe_id = $1 AND deleted_at IS NULL
	`

	queryInsertNewIngredient = `
		INSERT INTO ingredients
        (
			unit,
			cost, 
			recipe_id,
			ingredient_stock_id, 
			required_stock
		) VALUES ($1, $2, $3, $4, $5)
	`

	queryFindIngredientByRecipeID = `
		SELECT 
			id, 
			recipe_id, 
			ingredient_stock_id,
            unit, 
			cost, 
			required_stock
		FROM ingredients
		WHERE recipe_id = $1
	`
)
