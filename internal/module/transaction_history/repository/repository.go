package repository

import (
	"context"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_history/ports"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.TransactionHistoryRepository = &transactionHistoryRepository{}

type transactionHistoryRepository struct {
	db *sqlx.DB
}

func NewTransactionHistoryRepository(db *sqlx.DB) *transactionHistoryRepository {
	return &transactionHistoryRepository{
		db: db,
	}
}

func (r *transactionHistoryRepository) SumTotalAmount(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var sum int

	err := r.db.GetContext(ctx, &sum, r.db.Rebind(querySumTotalAmount), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Any("payload", startDate).Any("payload", endDate).Msg("repository::SumTotalAmount - Failed to sum total amount transaction history")
		return 0, err
	}

	return sum, err
}

func (r *transactionHistoryRepository) CountTransactionHistory(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var sum int

	err := r.db.GetContext(ctx, &sum, r.db.Rebind(queryCountTransactionHistory), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Any("payload", startDate).Any("payload", endDate).Msg("repository::CountTransactionHistory - Failed to count transaction history")
		return 0, err
	}

	return sum, err
}

func (r *transactionHistoryRepository) SumTotalQuantity(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var sum int

	err := r.db.GetContext(ctx, &sum, r.db.Rebind(querySumTotalQuantity), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Any("payload", startDate).Any("payload", endDate).Msg("repository::SumTotalQuantity - Failed to sum total quantity transaction history")
		return 0, err
	}

	return sum, err
}

func (r *transactionHistoryRepository) SumTotalCapital(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var sum int

	err := r.db.GetContext(ctx, &sum, r.db.Rebind(querySumTotalCapital), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Any("payload", startDate).Any("payload", endDate).Msg("repository::SumTotalCapital - Failed to sum total capital transaction history")
		return 0, err
	}

	return sum, err
}
