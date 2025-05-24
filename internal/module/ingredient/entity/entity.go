package entity

type Ingredient struct {
	ID                int     `db:"id"`
	Unit              string  `db:"unit"`
	Cost              float64 `db:"cost"`
	RecipeID          int     `db:"recipe_id"`
	IngredientStockID int     `db:"ingredient_stock_id"`
	RequiredStock     float64 `db:"required_stock"`
}
