package service

import (
	expensesPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/ports"
	"github.com/jmoiron/sqlx"
)

var _ expensesPorts.ExpensesService = &expensesService{}

type expensesService struct {
	db                 *sqlx.DB
	expensesRepository expensesPorts.ExpensesRepository
}

func NewExpensesService(
	db *sqlx.DB,
	expensesRepository expensesPorts.ExpensesRepository,
) *expensesService {
	return &expensesService{
		db:                 db,
		expensesRepository: expensesRepository,
	}
}
