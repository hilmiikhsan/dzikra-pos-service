package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.RecipeRepository = &recipeRepository{}

type recipeRepository struct {
	db *sqlx.DB
}

func NewRecipeRepository(db *sqlx.DB) *recipeRepository {
	return &recipeRepository{
		db: db,
	}
}

func (r *recipeRepository) InsertNewRecipe(ctx context.Context, tx *sqlx.Tx, data *entity.Recipe) (*entity.Recipe, error) {
	var res = new(entity.Recipe)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewRecipe),
		data.CapitalPrice,
	).Scan(
		&res.ID,
		&res.CapitalPrice,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewRecipe - Failed to insert new recipe")
		return nil, err
	}

	return res, nil
}
