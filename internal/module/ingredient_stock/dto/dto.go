package dto

type CreateIngredientStockRequest struct {
	Name                string `json:"name" validate:"required,max=100,min=1,xss_safe"`
	RequiredStock       int    `json:"required_stock" validate:"required,numeric,number"`
	Unit                string `json:"unit" validate:"required,max=10,min=1,xss_safe"`
	PricePerAmountStock int    `json:"price_per_amount_stock" validate:"required,numeric,number"`
	AmountStockPerPrice int    `json:"amount_stock_per_price" validate:"required,numeric,number"`
}

type CreateIngredientStockResponse struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	RequiredStock       int    `json:"required_stock"`
	Unit                string `json:"unit"`
	PricePerAmountStock int    `json:"price_per_amount_stock"`
	AmountStockPerPrice int    `json:"amount_stock_per_price"`
}

type GetListIngredientStockResponse struct {
	Stock       []GetListIngredientStock `json:"stock"`
	TotalPages  int                      `json:"total_page"`
	CurrentPage int                      `json:"current_page"`
	PageSize    int                      `json:"page_size"`
	TotalData   int                      `json:"total_data"`
}

type GetListIngredientStock struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	RequiredStock       int    `json:"required_stock"`
	Unit                string `json:"unit"`
	PricePerAmountStock int    `json:"price_per_amount_stock"`
	AmountStockPerPrice int    `json:"amount_stock_per_price"`
}
