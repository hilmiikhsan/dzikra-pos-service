package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/entity"
	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	InsertNewTransaction(ctx context.Context, tx *sqlx.Tx, data *entity.Transaction) error
	UpdateCashField(ctx context.Context, tx *sqlx.Tx, id, totalMoney, changeMoney string) error
	UpdateTransactionByID(ctx context.Context, tx *sqlx.Tx, data *entity.Transaction) error
	FindListTransaction(ctx context.Context, limit, offset int, search string) ([]dto.GetListTransaction, int, error)
	FindTransactionWithItemsByID(ctx context.Context, id string) (*entity.Transaction, error)
}

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *dto.CreateTransactionRequest, tableNumber int) (*dto.CreateTransactionResponse, error)
	GetListTransaction(ctx context.Context, page, limit int, search string) (*dto.GetListTransactionResponse, error)
	GetTransactionDetail(ctx context.Context, id string) (*dto.GetTransactionDetailResponse, error)
}
