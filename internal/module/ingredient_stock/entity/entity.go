package entity

type IngredientStock struct {
	ID                       int    `db:"id"`
	Name                     string `db:"name"`
	Unit                     string `db:"unit"`
	PricePerAmountStock      int    `db:"price_per_amount_stock"`
	RequiredStock            int    `db:"required_stock"`
	AmountPriceRequiredStock int    `db:"amount_price_required_stock"`
	AmountStockPerPrice      int    `db:"amount_stock_per_price"`
}
