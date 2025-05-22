package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductCategoryRepository = &productCategoryRepository{}

type productCategoryRepository struct {
	db *sqlx.DB
}

func NewProductCategoryRepository(db *sqlx.DB) *productCategoryRepository {
	return &productCategoryRepository{
		db: db,
	}
}

func (r *productCategoryRepository) InsertNewProductCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ProductCategory) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProductCategory),
		data.Name,
	).Scan(
		&res.ID,
		&res.Name,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"product_categories_name_key": constants.ErrProductCategoryAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewProductCategory - Failed to insert new product category")
			return nil, handleErr
		}

		if productCategory, ok := val.(*entity.ProductCategory); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewProductCategory - Failed to insert new product category")
			return productCategory, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewProductCategory - Failed to insert new product category")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) UpdateProductCategory(ctx context.Context, tx *sqlx.Tx, data *entity.ProductCategory) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProductCategory),
		data.Name,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateProductCategory - product category with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductCategoryNotFound)
		}

		uniqueConstraints := map[string]string{
			"product_categories_name_key": constants.ErrProductCategoryAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::UpdateProductCategory - Failed to update new product category")
			return nil, handleErr
		}

		if productCategory, ok := val.(*entity.ProductCategory); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::UpdateProductCategory - Failed to update new product category")
			return productCategory, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateProductCategory - Failed to update product category")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) FindProductCategoryByID(ctx context.Context, id int) (*entity.ProductCategory, error) {
	var res = new(entity.ProductCategory)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindProductCategoryByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindProductCategoryByID - product category not found")
			return nil, errors.New(constants.ErrProductCategoryNotFound)
		}

		log.Error().Err(err).Msg("repository::FindProductCategoryByID - failed to get product category")
		return nil, err
	}

	return res, nil
}

func (r *productCategoryRepository) FindListProductCategory(ctx context.Context, limit, offset int, search string) ([]dto.GetProductCategoryResponse, int, error) {
	args := []interface{}{search, limit, offset}

	var responses []entity.ProductCategory
	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListproductCategory), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListproductCategory - error executing query")
		return nil, 0, err
	}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListproductCategory), args[:1]...); err != nil {
		log.Error().Err(err).Msg("repository::FindListproductCategory - error counting product categories")
		return nil, 0, err
	}

	productCategories := make([]dto.GetProductCategoryResponse, 0, len(responses))
	for _, v := range responses {
		productCategories = append(productCategories, dto.GetProductCategoryResponse{
			ID:   v.ID,
			Name: v.Name,
		})
	}

	return productCategories, total, nil
}

func (r *productCategoryRepository) SoftDeleteProductCategoryByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, tx.Rebind(querySoftDeleteProductCategoryByID), id)
	if err != nil {
		log.Error().Err(err).Any("id", id).Msg("repository::SoftDeleteProductCategoryByID - Failed to soft delete product category")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteProductCategoryByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrProductCategoryNotFound)
		log.Error().Err(errNotFound).Any("id", id).Msg("repository::SoftDeleteProductCategoryByID - product Category not found")
		return errNotFound
	}

	return nil
}
