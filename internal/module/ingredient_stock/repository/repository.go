package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/ports"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var _ ports.IngredientStockRepository = &ingredientStockRepository{}

type ingredientStockRepository struct {
	db *sqlx.DB
}

func NewIngredientStockRepository(db *sqlx.DB) *ingredientStockRepository {
	return &ingredientStockRepository{
		db: db,
	}
}

func (r *ingredientStockRepository) InsertNewIngredientStock(ctx context.Context, tx *sqlx.Tx, data *entity.IngredientStock) (*entity.IngredientStock, error) {
	var res = new(entity.IngredientStock)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewIngredientStock),
		data.Name,
		data.Unit,
		data.PricePerAmountStock,
		data.RequiredStock,
		data.AmountPriceRequiredStock,
		data.AmountStockPerPrice,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Unit,
		&res.PricePerAmountStock,
		&res.RequiredStock,
		&res.AmountPriceRequiredStock,
		&res.AmountStockPerPrice,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewIngredientStock - Failed to insert new ingredient stock")
		return nil, err
	}

	return res, nil
}

func (r *ingredientStockRepository) UpdateIngredientStock(ctx context.Context, tx *sqlx.Tx, data *entity.IngredientStock) (*entity.IngredientStock, error) {
	var res = new(entity.IngredientStock)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateIngredientStock),
		data.Name,
		data.Unit,
		data.PricePerAmountStock,
		data.RequiredStock,
		data.AmountPriceRequiredStock,
		data.AmountStockPerPrice,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Unit,
		&res.PricePerAmountStock,
		&res.RequiredStock,
		&res.AmountPriceRequiredStock,
		&res.AmountStockPerPrice,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateIngredientStock - member with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrIngredientStockNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateIngredientStock - Failed to update ingredient stock")
		return nil, err
	}

	return res, nil
}

func (r *ingredientStockRepository) FindListIngredientStock(ctx context.Context, limit, offset int, search string) ([]dto.GetListIngredientStock, int, error) {
	args := []interface{}{search, limit, offset}

	var responses []entity.IngredientStock
	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListIngredientStock), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListIngredientStock - error executing query")
		return nil, 0, err
	}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListIngredientStock), args[:1]...); err != nil {
		log.Error().Err(err).Msg("repository::FindListIngredientStock - error counting ingredient stock")
		return nil, 0, err
	}

	ingredientStocks := make([]dto.GetListIngredientStock, 0, len(responses))
	for _, v := range responses {
		ingredientStocks = append(ingredientStocks, dto.GetListIngredientStock{
			ID:                  v.ID,
			Name:                v.Name,
			Unit:                v.Unit,
			PricePerAmountStock: v.PricePerAmountStock,
			RequiredStock:       v.RequiredStock,
			AmountStockPerPrice: v.AmountStockPerPrice,
		})
	}

	return ingredientStocks, total, nil
}

func (r *ingredientStockRepository) FindIngredientStockByID(ctx context.Context, id int) (*entity.IngredientStock, error) {
	var res = new(entity.IngredientStock)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindIngredientStockByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Err(err).Msg("repository::FindIngredientStockByID - No ingredient stock found")
			return nil, errors.New(constants.ErrIngredientStockNotFound)
		}

		log.Error().Err(err).Msg("repository::FindIngredientStockByID - Failed to find ingredient stock")
		return nil, err
	}

	return res, nil
}

func (r *ingredientStockRepository) SoftDeleteIngredientStockByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, tx.Rebind(querySoftDeleteIngredientStockByID), id)
	if err != nil {
		log.Error().Err(err).Any("id", id).Msg("repository::SoftDeleteIngredientStockByID - Failed to soft delete ingredient stock")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteIngredientStockByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrIngredientStockNotFound)
		log.Error().Err(errNotFound).Any("id", id).Msg("repository::SoftDeleteIngredientStockByID - Ingredient stock not found")
		return errNotFound
	}

	return nil
}

func (r *ingredientStockRepository) FindIngredientStockByIDs(ctx context.Context, ids []int) ([]entity.IngredientStock, error) {
	var res []entity.IngredientStock

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindIngredientStockByIDs), pq.Array(ids))
	if err != nil {
		log.Error().Err(err).Msg("repository::FindIngredientStockByIDs - Failed to get ingredient stock by ids")
		return nil, err
	}

	return res, err
}

func (r *ingredientStockRepository) CountIngredientStockByID(ctx context.Context, id int) (int, error) {
	var count int

	err := r.db.GetContext(ctx, &count, r.db.Rebind(queryCountIngredientStockByID), id)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountIngredientStockByID - Failed to count ingredient stock")
		return 0, err
	}

	return count, nil
}
