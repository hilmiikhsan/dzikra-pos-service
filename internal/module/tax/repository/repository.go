package repository

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.TaxRepository = &taxRepository{}

type taxRepository struct {
	db *sqlx.DB
}

func NewTaxRepository(db *sqlx.DB) *taxRepository {
	return &taxRepository{
		db: db,
	}
}

func (r *taxRepository) InsertNewTax(ctx context.Context, tx *sqlx.Tx, data *entity.Tax) (*entity.Tax, error) {
	var res = new(entity.Tax)

	err := tx.QueryRowContext(ctx, tx.Rebind(queryInsertNewTax),
		data.Tax,
	).Scan(
		&res.ID,
		&res.Tax,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewTax - Failed to insert new tax")
		return nil, err
	}

	return res, nil
}

func (r *taxRepository) FindTax(ctx context.Context) (*entity.Tax, error) {
	var res = new(entity.Tax)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindTax))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Err(err).Msg("repository::FindTax - No tax found")
			return nil, nil
		}

		log.Error().Err(err).Msg("repository::FindTax - Failed to find tax")
		return nil, err
	}

	return res, nil
}

func (r *taxRepository) UpdateTax(ctx context.Context, tx *sqlx.Tx, data *entity.Tax) (*entity.Tax, error) {
	var res = new(entity.Tax)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateTax),
		data.Tax,
		data.ID,
	).Scan(
		&res.ID,
		&res.Tax,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateTax - Failed to update tax")
		return nil, err
	}

	return res, nil
}
