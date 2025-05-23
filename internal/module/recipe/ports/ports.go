package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	"github.com/jmoiron/sqlx"
)

type RecipeRepository interface {
	InsertNewRecipe(ctx context.Context, tx *sqlx.Tx, data *entity.Recipe) (*entity.Recipe, error)
}
