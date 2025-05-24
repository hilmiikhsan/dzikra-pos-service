package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient/ports"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.IngredientRepository = &ingredientRepository{}

type ingredientRepository struct {
	db *sqlx.DB
}

func NewIngredientRepository(db *sqlx.DB) *ingredientRepository {
	return &ingredientRepository{
		db: db,
	}
}

func (r *ingredientRepository) FindIngredientByRecipeIDs(ctx context.Context, recipeIDs []int) ([]entity.Ingredient, error) {
	var res []entity.Ingredient

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindIngredientByRecipeIDs), pq.Array(recipeIDs))
	if err != nil {
		log.Error().Err(err).Any("payload", recipeIDs).Msg("repository::FindIngredientByRecipeIDs - Failed to get ingredient by recipe ids")
		return nil, err
	}

	return res, err
}

func (r *ingredientRepository) DeleteIngredientByRecipeID(ctx context.Context, tx *sqlx.Tx, recipeID int) error {
	if _, err := tx.ExecContext(ctx, r.db.Rebind(queryDeleteIngredientByRecipeID), recipeID); err != nil {
		log.Error().Err(err).Int("recipe_id", recipeID).Msg("repository::DeleteIngredientByRecipeID - failed to delete ingredient by recipe id")
		return err
	}

	return nil
}

func (r *ingredientRepository) InsertNewIngredient(ctx context.Context, tx *sqlx.Tx, data *entity.Ingredient) error {
	if _, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewIngredient),
		data.Unit,
		data.Cost,
		data.RecipeID,
		data.IngredientStockID,
		data.RequiredStock,
	); err != nil {
		log.Error().Err(err).Interface("data", data).Msg("repository::InsertNewIngredient - failed to insert new ingredient")
		return err
	}

	return nil
}

func (r *ingredientRepository) FindIngredientByRecipeID(ctx context.Context, tx *sqlx.Tx, recipeID int) ([]entity.Ingredient, error) {
	var res []entity.Ingredient

	if err := tx.SelectContext(ctx, &res, r.db.Rebind(queryFindIngredientByRecipeID), recipeID); err != nil {
		log.Error().Err(err).Int("recipe_id", recipeID).Msg("repository::FindIngredientByRecipeID - failed to find ingredient by recipe id")
		return nil, err
	}

	return res, nil
}
