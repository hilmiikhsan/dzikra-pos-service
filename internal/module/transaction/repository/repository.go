package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/ports"
	transactionItemEntity "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/entity"
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

func (r *transactionRepository) FindListTransaction(ctx context.Context, limit, offset int, search string) ([]dto.GetListTransaction, int, error) {
	args := []interface{}{search, search, search, limit, offset}

	var ents []entity.Transaction
	if err := r.db.SelectContext(ctx, &ents, r.db.Rebind(queryFindListTransaction), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListTransaction - error executing query")
		return nil, 0, err
	}

	countArgs := []interface{}{search, search, search}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListTransaction), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListTransaction - error counting transactions")
		return nil, 0, err
	}

	out := make([]dto.GetListTransaction, 0, len(ents))
	for _, v := range ents {
		out = append(out, dto.GetListTransaction{
			ID:                       v.ID.String(),
			Status:                   v.Status,
			PhonenUmber:              v.PhoneNumber,
			Name:                     v.Name,
			Email:                    v.Email,
			IsMember:                 v.IsMember,
			TotalQuantity:            fmt.Sprintf("%d", v.TotalQuantity),
			TotalProductAmount:       fmt.Sprintf("%d", v.TotalProductAmount),
			TotalProductCapitalPrice: fmt.Sprintf("%d", v.TotalProductCapitalPrice),
			TotalAmount:              fmt.Sprintf("%d", v.TotalAmount),
			DiscountPercentage:       fmt.Sprintf("%d", v.DiscountPercentage),
			VTransactionID:           v.VTransactionID,
			VPaymentID:               v.VPaymentID,
			VPaymentRedirectUrl:      v.VPaymentRedirectUrl,
			PaymentType:              v.PaymentType,
			TableNumber:              fmt.Sprintf("%d", v.TableNumber),
			CreatedAt:                v.CreatedAt.Format(time.RFC3339),
			Notes:                    v.Notes,
			TaxAmount:                fmt.Sprintf("%d", v.TaxAmount),
		})
	}

	return out, total, nil
}

func (r *transactionRepository) FindTransactionWithItemsByID(ctx context.Context, id string) (*entity.Transaction, error) {
	var txRow entity.Transaction

	if err := r.db.GetContext(ctx, &txRow, r.db.Rebind(queryFindTransactionWithItemsByID), id); err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("id", id).Msg("repository::FindTransactionWithItemsByID - transaction not found")
			return nil, errors.New(constants.ErrTransactionNotFound)
		}

		log.Error().Err(err).Msg("repository::FindTransactionWithItemsByID - failed to get transaction")
		return nil, err
	}

	var items []transactionItemEntity.TransactionItem
	if err := r.db.SelectContext(ctx, &items, r.db.Rebind(queryFindItemsByTransactionID), id); err != nil {
		if err == sql.ErrNoRows {
			log.Warn().Str("id", id).Msg("repository::FindTransactionWithItemsByID - transaction not found")
			return nil, errors.New(constants.ErrTransactionNotFound)
		}

		log.Error().Err(err).Msg("repository::FindTransactionWithItemsByID - failed to get transaction items")
		return nil, err
	}

	txRow.TransactionItems = items
	return &txRow, nil
}
