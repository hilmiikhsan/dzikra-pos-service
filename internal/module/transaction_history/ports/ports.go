package ports

import (
	"context"
	"time"
)

type TransactionHistoryRepository interface {
	SumTotalAmount(ctx context.Context, startDate, endDate time.Time) (int, error)
	CountTransactionHistory(ctx context.Context, startDate, endDate time.Time) (int, error)
	SumTotalQuantity(ctx context.Context, startDate, endDate time.Time) (int, error)
	SumTotalCapital(ctx context.Context, startDate, endDate time.Time) (int, error)
}
