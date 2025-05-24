package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/ports"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func (r *recipeRepository) CountRecipe(ctx context.Context) (int, error) {
	var count int

	err := r.db.GetContext(ctx, &count, r.db.Rebind(queryCountRecipe))
	if err != nil {
		log.Error().Err(err).Msg("repository::CountRecipe - Failed to count recipe")
		return 0, err
	}

	return count, nil
}

func (r *recipeRepository) FindListRecipeByProductIDs(ctx context.Context, productIDs []int) ([]entity.Recipe, error) {
	var res []entity.Recipe

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindListRecipeByProductIDs), productIDs)
	if err != nil {
		log.Error().Err(err).Msg("repository::FindListRecipeByProductIDs - Failed to count recipe")
		return nil, err
	}

	return res, nil
}

func (r *recipeRepository) FindRecipeByIDs(ctx context.Context, ids []int) ([]entity.Recipe, error) {
	var recs []entity.Recipe

	if err := r.db.SelectContext(ctx, &recs, r.db.Rebind(queryFindRecipesByIDs), pq.Array(ids)); err != nil {
		log.Error().Err(err).Msg("repository::FindRecipeByIDs - Failed to find recipe by ids")
		return nil, err
	}

	return recs, nil
}

func (r *recipeRepository) UpdateRecipeCapitalPrice(ctx context.Context, tx *sqlx.Tx, data *entity.Recipe) error {
	if _, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateRecipeCapitalPrice), data.ID, data.CapitalPrice); err != nil {
		log.Error().Err(err).Int("recipe_id", data.ID).Int("capPrice", data.CapitalPrice).Msg("repository::UpdateRecipeCapitalPrice - failed to update recipe capital price")
		return err
	}

	return nil
}

func (r *recipeRepository) FindDetailRecipes(ctx context.Context, recipeID int) (*entity.RecipeDetail, error) {
	var detail entity.RecipeDetail
	err := r.db.QueryRowxContext(ctx, r.db.Rebind(queryFindDetailRecipes), recipeID).Scan(
		&detail.ID,
		&detail.CapitalPrice,
		&detail.Product.ID,
		&detail.Product.Name,
		&detail.Product.Stock,
	)
	if err != nil {
		log.Error().Err(err).Int("recipe_id", recipeID).Msg("repository::FindDetailRecipes - failed to scan header")
		return nil, err
	}

	type ingRow struct {
		ID            int    `db:"id"`
		Unit          string `db:"unit"`
		Cost          int    `db:"cost"`
		RequiredStock int    `db:"required_stock"`
		RecipeID      int    `db:"recipe_id"`
		StockID       int    `db:"stock_id"`
		StockName     string `db:"stock_name"`
	}

	var rows []ingRow
	if err := r.db.SelectContext(ctx, &rows, r.db.Rebind(queryFindDetailRecipesWIthIngredients), recipeID); err != nil {
		log.Error().Err(err).Int("recipe_id", recipeID).Msg("repository::FindDetailRecipes - failed to fetch ingredients")
		return nil, err
	}

	detail.Ingredients = make([]entity.IngredientDetail, len(rows))
	for i, row := range rows {
		detail.Ingredients[i] = entity.IngredientDetail{
			ID:            row.ID,
			Unit:          row.Unit,
			Cost:          row.Cost,
			RequiredStock: row.RequiredStock,
			RecipeID:      row.RecipeID,
			Stock: entity.IngredientStockEntity{
				ID:   row.StockID,
				Name: row.StockName,
			},
		}
	}

	return &detail, nil
}
