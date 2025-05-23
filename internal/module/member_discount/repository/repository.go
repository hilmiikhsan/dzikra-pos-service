package repository

import (
	"context"
	"database/sql"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.MemberDiscountRepository = &memberDiscountRepository{}

type memberDiscountRepository struct {
	db *sqlx.DB
}

func NewMemberDiscountRepository(db *sqlx.DB) *memberDiscountRepository {
	return &memberDiscountRepository{
		db: db,
	}
}

func (r *memberDiscountRepository) InsertNewMemberDiscount(ctx context.Context, tx *sqlx.Tx, data *entity.MemberDiscount) (*entity.MemberDiscount, error) {
	var res = new(entity.MemberDiscount)

	err := tx.QueryRowContext(ctx, tx.Rebind(queryInsertNewMemberDiscount),
		data.Discount,
	).Scan(
		&res.ID,
		&res.Discount,
		&res.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewMember - Failed to insert new member discount")
		return nil, err
	}

	return res, nil
}

func (r *memberDiscountRepository) FindDiscount(ctx context.Context) (*entity.MemberDiscount, error) {
	var res = new(entity.MemberDiscount)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindMemberDiscount))
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Err(err).Msg("repository::FindDiscount - No member discount found")
			return nil, nil
		}

		log.Error().Err(err).Msg("repository::FindDiscount - Failed to find member discount")
		return nil, err
	}

	return res, nil
}

func (r *memberDiscountRepository) UpdateMemberDiscount(ctx context.Context, tx *sqlx.Tx, data *entity.MemberDiscount) (*entity.MemberDiscount, error) {
	var res = new(entity.MemberDiscount)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateMemberDiscount),
		data.Discount,
		data.ID,
	).Scan(
		&res.ID,
		&res.Discount,
		&res.UpdatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateMemberDiscount - Failed to update member discount")
		return nil, err
	}

	return res, nil
}
