package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/entity"
	"github.com/jmoiron/sqlx"
)

type TransactionItemRepository interface {
	InsertNewTransactionItem(ctx context.Context, tx *sqlx.Tx, data *entity.TransactionItem) error
	FindTransactionItemByTransactionID(ctx context.Context, transactionID string) ([]*entity.TransactionItem, error)
}
