package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/ports"
	product_category "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ProductRepository = &productRepository{}

type productRepository struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) *productRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) InsertNewProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error) {
	var res = new(entity.Product)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewProduct),
		data.Name,
		data.Description,
		data.ImageUrl,
		data.Stock,
		data.IsActive,
		data.ProductCategoryID,
		data.RealPrice,
		data.RecipeID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Description,
		&res.ImageUrl,
		&res.Stock,
		&res.IsActive,
		&res.ProductCategoryID,
		&res.RealPrice,
		&res.RecipeID,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewProduct - Failed to insert new product")
		return nil, err
	}

	return res, nil
}

func (r *productRepository) UpdateProduct(ctx context.Context, tx *sqlx.Tx, data *entity.Product) (*entity.Product, error) {
	var res = new(entity.Product)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateProduct),
		data.Name,
		data.Description,
		data.ImageUrl,
		data.Stock,
		data.IsActive,
		data.ProductCategoryID,
		data.RealPrice,
		data.RecipeID,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Description,
		&res.ImageUrl,
		&res.Stock,
		&res.IsActive,
		&res.ProductCategoryID,
		&res.RealPrice,
		&res.RecipeID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateProduct - product with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateProduct - Failed to update product")
		return nil, err
	}

	return res, nil
}

func (r *productRepository) CountProductByID(ctx context.Context, id int) (int, error) {
	var count int

	err := r.db.GetContext(ctx, &count, r.db.Rebind(queryCountProductByID), id)
	if err != nil {
		log.Error().Err(err).Msg("repository::CountProductByID - Failed to count product by id")
		return 0, err
	}

	return count, nil
}

func (r *productRepository) FindProductByID(ctx context.Context, id int) (*entity.Product, error) {
	var res = new(entity.Product)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindProductByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::FindProductByID - product with id %d is not found", id)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrProductNotFound)
		}

		log.Error().Err(err).Msg("repository::FindProductByID - Failed to find product by id")
		return nil, err
	}

	return res, nil
}

func (r *productRepository) FindListProduct(ctx context.Context, limit, offset int, search string) ([]dto.GetListProduct, int, error) {
	args := []interface{}{
		search, search, limit, offset,
	}

	var ents []entity.Product
	if err := r.db.SelectContext(ctx, &ents, r.db.Rebind(queryFindAllProduct), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindAllProduct - error executing query")
		return nil, 0, err
	}

	countArgs := []interface{}{search, search}
	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindAllProduct), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindAllProduct - error counting articles")
		return nil, 0, err
	}

	publicURL := config.Envs.MinioStorage.PublicURL
	out := make([]dto.GetListProduct, 0, len(ents))
	for _, v := range ents {
		out = append(out, dto.GetListProduct{
			ID:          v.ID,
			Name:        v.Name,
			Description: v.Description,
			Stock:       v.Stock,
			RealPrice:   v.RealPrice,
			IsActive:    v.IsActive,
			ImageUrl:    utils.FormatMediaPathURL(v.ImageUrl, publicURL),
			ProductCategory: product_category.GetProductCategoryResponse{
				ID:   v.ProductCategoryID,
				Name: v.ProductCategoryName,
			},
		})
	}

	return out, total, nil
}

func (r *productRepository) SoftDeleteProductByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, r.db.Rebind(qyerySoftDeleteArticleByID), id)
	if err != nil {
		log.Error().Err(err).Int("id", id).Msg("repository::SoftDeleteProductByID - Failed to soft delete product")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteProductByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrProductNotFound)
		log.Error().Err(errNotFound).Int("id", id).Msg("repository::SoftDeleteProductByID - product not found")
		return errNotFound
	}

	return nil
}
