package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.TransactionRepository = &transactionRepository{}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) InsertNewTransaction(ctx context.Context, tx *sqlx.Tx, data *entity.Transaction) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewTransaction),
		data.ID,
		data.Status,
		data.PhoneNumber,
		data.Name,
		data.Email,
		data.IsMember,
		data.TotalQuantity,
		data.TotalProductAmount,
		data.TotalAmount,
		data.VPaymentID,
		data.VPaymentRedirectUrl,
		data.VTransactionID,
		data.DiscountPercentage,
		data.ChangeMoney,
		data.PaymentType,
		data.TotalMoney,
		data.TableNumber,
		data.TotalProductCapitalPrice,
		data.TaxAmount,
		data.Notes,
		data.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewTransaction - Failed to insert new transaction")
		return err
	}

	return nil
}

func (r *transactionRepository) UpdateCashField(ctx context.Context, tx *sqlx.Tx, id, totalMoney, changeMoney string) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateCashField),
		totalMoney,
		changeMoney,
		id,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", id).Msg("repository::UpdateCashField - Failed to update cash field transaction")
		return err
	}

	return nil
}

func (r *transactionRepository) UpdateTransactionByID(ctx context.Context, tx *sqlx.Tx, data *entity.Transaction) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryUpdateTransactionByID),
		data.VTransactionID,
		data.VPaymentID,
		data.VPaymentRedirectUrl,
		data.ID,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateTransactionByID - Failed to update transaction by ID")
		return err
	}

	return nil
}
