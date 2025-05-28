package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.ExpensesRepository = &expensesRepository{}

type expensesRepository struct {
	db *sqlx.DB
}

func NewExpensesRepository(db *sqlx.DB) *expensesRepository {
	return &expensesRepository{
		db: db,
	}
}

func (r *expensesRepository) InsertNewExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error) {
	var res = new(entity.Expenses)

	err := tx.QueryRowContext(ctx, tx.Rebind(queryInsertNewExpenses),
		data.Name,
		data.Cost,
		data.Date,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Cost,
		&res.Date,
		&res.CreatedAt,
	)
	if err != nil {
		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewExpenses - Failed to insert new expenses")
		return nil, err
	}

	return res, nil
}

func (r *expensesRepository) FindListExpenses(ctx context.Context, limit, offset int, search string) ([]dto.GetListExpenses, int, error) {
	args := []interface{}{search, search, limit, offset}

	var responses []entity.Expenses
	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListExpenses), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListExpenses - error executing query")
		return nil, 0, err
	}

	countArgs := []interface{}{search, search}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListExpenses), countArgs...); err != nil {
		log.Error().Err(err).Msg("repository::FindListExpenses - error counting total expenses")
		return nil, 0, err
	}

	members := make([]dto.GetListExpenses, 0, len(responses))
	for _, v := range responses {
		members = append(members, dto.GetListExpenses{
			ID:        v.ID,
			Name:      v.Name,
			Cost:      v.Cost,
			Date:      utils.FormatTime(v.Date),
			CreatedAt: utils.FormatTime(v.CreatedAt),
			UpdatedAt: utils.FormatTime(v.UpdatedAt),
		})
	}

	return members, total, nil
}

func (r *expensesRepository) FindExpensesByID(ctx context.Context, id int) (*entity.Expenses, error) {
	var res = new(entity.Expenses)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindExpensesByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Err(err).Any("id", id).Msg("repository::FindExpensesByID - No expenses found with given ID")
			return nil, errors.New(constants.ErrExpensesNotFound)
		}

		log.Error().Err(err).Any("id", id).Msg("repository::FindExpensesByID - Failed to find expenses by ID")
		return nil, err
	}

	return res, nil
}

func (r *expensesRepository) UpdateExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error) {
	var res = new(entity.Expenses)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateExpenses),
		data.Name,
		data.Cost,
		data.Date,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.Cost,
		&res.Date,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Debug().Err(err).Any("id", data.ID).Msg("repository::UpdateExpenses - No expenses found with given ID")
			return nil, errors.New(constants.ErrExpensesNotFound)
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateExpenses - Failed to update expenses")
		return nil, err
	}

	return res, nil
}

func (r *expensesRepository) SoftDeleteExpensesByID(ctx context.Context, tx *sqlx.Tx, id int) error {
	result, err := tx.ExecContext(ctx, tx.Rebind(querySoftDeleteExpensesByID), id)
	if err != nil {
		log.Error().Err(err).Any("id", id).Msg("repository::SoftDeleteExpensesByID - Failed to soft delete expenses by ID")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteExpensesByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrExpensesNotFound)
		log.Error().Err(errNotFound).Any("id", id).Msg("repository::SoftDeleteExpensesByID - MeExpensesmber not found")
		return errNotFound
	}

	return nil
}

func (r *expensesRepository) SumTotalCostExpenses(ctx context.Context, startDate, endDate time.Time) (int, error) {
	var sum int

	err := r.db.GetContext(ctx, &sum, r.db.Rebind(querySumTotalCostExpenses), startDate, endDate)
	if err != nil {
		log.Error().Err(err).Any("payload", startDate).Any("payload", endDate).Msg("repository::SumTotalCostExpenses - Failed to sum total amount expenses")
		return 0, err
	}

	return sum, err
}
