package dto

type GetListRecipeResponse struct {
	RecipesSerialize []ProductWithRecipe `json:"recipesSerialize"`
	TotalPages       int                 `json:"total_pages"`
	CurrentPage      int                 `json:"current_page"`
	PageSize         int                 `json:"page_size"`
	TotalData        int                 `json:"total_data"`
}

type Ingredient struct {
	ID              int             `json:"id"`
	Unit            string          `json:"unit"`
	Cost            int             `json:"cost"`
	RequiredStock   int             `json:"required_stock"`
	RecipeID        int             `json:"recipe_id"`
	IngredientStock IngredientStock `json:"IngredientStock"`
}

type IngredientStock struct {
	ID                       int    `json:"id"`
	Name                     string `json:"name"`
	RequiredStock            int    `json:"required_stock"`
	Unit                     string `json:"unit"`
	PricePerAmountStock      int    `json:"price_per_amount_stock"`
	AmountStockPerPrice      int    `json:"amount_stock_per_price"`
	AmountPriceRequiredStock int    `json:"amount_price_required_stock"`
}

type Recipe struct {
	ID           int          `json:"id"`
	CapitalPrice int          `json:"capital_price"`
	Ingredients  []Ingredient `json:"Ingredients"`
}

type ProductWithRecipe struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Recipes Recipe `json:"recipes"`
}

type UpdateRecipeRequest struct {
	Ingredients []IngredientRequest `json:"ingredients" validate:"required,non_empty_array"`
}

type IngredientRequest struct {
	IngredientID        string `json:"ingredient_id" validate:"required,min=1,max=20,xss_safe"`
	RequiredStock       int    `json:"required_stock" validate:"required,numeric,number,gt=1"`
	Unit                string `json:"unit" validate:"required,min=1,max=20,xss_safe"`
	Cost                int    `json:"cost" validate:"required,numeric,number,gt=0"`
	PricePerAmountStock int    `json:"price_per_amount_stock" validate:"required,numeric,number,gt=0"`
}

type UpdateRecipeResponse struct {
	ID           int                   `json:"id"`
	CapitalPrice int                   `json:"capital_price"`
	Product      ProductDTO            `json:"Product"`
	Ingredients  []IngredientDetailDTO `json:"Ingredients"`
}

type ProductDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
}

type IngredientDetailDTO struct {
	ID              int                `json:"id"`
	Unit            string             `json:"unit"`
	Cost            int                `json:"cost"`
	RequiredStock   string             `json:"required_stock"`
	RecipeID        int                `json:"recipe_id"`
	IngredientStock IngredientStockDTO `json:"IngredientStock"`
}

type IngredientStockDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
