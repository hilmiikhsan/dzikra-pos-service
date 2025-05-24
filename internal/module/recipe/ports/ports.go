package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	"github.com/jmoiron/sqlx"
)

type RecipeRepository interface {
	InsertNewRecipe(ctx context.Context, tx *sqlx.Tx, data *entity.Recipe) (*entity.Recipe, error)
	CountRecipe(ctx context.Context) (int, error)
	FindListRecipeByProductIDs(ctx context.Context, productIDs []int) ([]entity.Recipe, error)
	FindRecipeByIDs(ctx context.Context, ids []int) ([]entity.Recipe, error)
	UpdateRecipeCapitalPrice(ctx context.Context, tx *sqlx.Tx, data *entity.Recipe) error
	FindDetailRecipes(ctx context.Context, recipeID int) (*entity.RecipeDetail, error)
}

type RecipeService interface {
	GetListRecipe(ctx context.Context, page, limit int, search string) (*dto.GetListRecipeResponse, error)
	UpdateRecipe(ctx context.Context, req *dto.UpdateRecipeRequest, productID int) (*dto.UpdateRecipeResponse, error)
}
