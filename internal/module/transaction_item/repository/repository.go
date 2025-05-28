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

func (r *transactionItemRepository) FindTransactionItemByTransactionID(ctx context.Context, transactionID string) ([]*entity.TransactionItem, error) {
	var res []*entity.TransactionItem

	err := r.db.SelectContext(ctx, &res, r.db.Rebind(queryFindTransactionItemByTransactionID), transactionID)
	if err != nil {
		log.Error().Err(err).Str("transactionID", transactionID).Msg("repository::FindTransactionItemByTransactionID - Failed to find transaction items")
		return nil, err
	}

	return res, nil
}
