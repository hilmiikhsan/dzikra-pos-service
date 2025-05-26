package entity

type Recipe struct {
	ID           int `db:"id"`
	CapitalPrice int `db:"capital_price"`
}

type IngredientEntity struct {
	ID                  int
	RecipeID            int
	IngredientStockID   int
	Unit                string
	Cost                int
	RequiredStock       int
	PricePerAmountStock int
}

type IngredientStockEntity struct {
	ID            int
	RequiredStock int
	Name          string
}

type ProductEntity struct {
	ID       int
	Name     string
	Stock    int
	RecipeID int
}
type IngredientDetail struct {
	ID            int    `db:"id"`
	Unit          string `db:"unit"`
	Cost          int    `db:"cost"`
	RequiredStock int    `db:"required_stock"`
	RecipeID      int    `db:"recipe_id"`
	Stock         IngredientStockEntity
}

type RecipeDetail struct {
	ID           int `db:"id"`
	CapitalPrice int `db:"capital_price"`
	Product      ProductEntity
	Ingredients  []IngredientDetail
}

type IngredientInfo struct {
	RecipeID     string `db:"recipe_id"`
	IngredientID string `db:"ingredient_id"`
	ReqPerUnit   int    `db:"required_stock"`
	StockID      string `db:"stock_id"`
	StockAmount  int    `db:"stock_amount"`
}
