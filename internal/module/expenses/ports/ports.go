package ports

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/entity"
	"github.com/jmoiron/sqlx"
)

type ExpensesRepository interface {
	InsertNewExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error)
	FindListExpenses(ctx context.Context, limit, offset int, search string) ([]dto.GetListExpenses, int, error)
	FindExpensesByID(ctx context.Context, id int) (*entity.Expenses, error)
	UpdateExpenses(ctx context.Context, tx *sqlx.Tx, data *entity.Expenses) (*entity.Expenses, error)
	SoftDeleteExpensesByID(ctx context.Context, tx *sqlx.Tx, id int) error
}

type ExpensesService interface {
	CreateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest) (*dto.CreateOrUpdateExpensesResponse, error)
	GetListExpenses(ctx context.Context, page, limit int, search string) (*dto.GetListExpensesResponse, error)
	GetDetailExpenses(ctx context.Context, id int) (*dto.GetListExpenses, error)
	UpdateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest, id int) (*dto.CreateOrUpdateExpensesResponse, error)
	RemoveExpenses(ctx context.Context, id int) error
}
