package repository

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.TransactionItemRepository = &transactionItemRepository{}

type transactionItemRepository struct {
	db *sqlx.DB
}

func NewTransactionItemRepository(db *sqlx.DB) *transactionItemRepository {
	return &transactionItemRepository{
		db: db,
	}
}

func (r *transactionItemRepository) InsertNewTransactionItem(ctx context.Context, tx *sqlx.Tx, data *entity.TransactionItem) error {
	_, err := tx.ExecContext(ctx, r.db.Rebind(queryInsertNewTransactionItem),
		data.ID,
		data.Quantity,
		data.TotalAmount,
		data.ProductName,
		data.ProductPrice,
		data.TransactionID,
		data.ProductID,
		data.ProductCapitalPrice,
		data.TotalAmountCapitalPrice,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewTransactionItem - Failed to insert new transaction item")
		return err
	}

	return nil
}
