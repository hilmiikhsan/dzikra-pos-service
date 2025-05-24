package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/entity"
	"github.com/jmoiron/sqlx"
)

type IngredientRepository interface {
	FindIngredientByRecipeIDs(ctx context.Context, recipeIDs []int) ([]entity.Ingredient, error)
	DeleteIngredientByRecipeID(ctx context.Context, tx *sqlx.Tx, recipeID int) error
	InsertNewIngredient(ctx context.Context, tx *sqlx.Tx, data *entity.Ingredient) error
	FindIngredientByRecipeID(ctx context.Context, tx *sqlx.Tx, recipeID int) ([]entity.Ingredient, error)
}
