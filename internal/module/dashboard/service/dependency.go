package service

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/dashboard/ports"
	expensesPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/ports"
	transactionHistoryPorts "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_history/ports"
	"github.com/jmoiron/sqlx"
)

var _ ports.DashboardService = &dashboardService{}

type dashboardService struct {
	db                           *sqlx.DB
	transactionHistoryRepository transactionHistoryPorts.TransactionHistoryRepository
	expensesRepository           expensesPorts.ExpensesRepository
}

func NewDashboardService(
	db *sqlx.DB,
	transactionHistoryRepository transactionHistoryPorts.TransactionHistoryRepository,
	expensesRepository expensesPorts.ExpensesRepository,
) *dashboardService {
	return &dashboardService{
		db:                           db,
		transactionHistoryRepository: transactionHistoryRepository,
		expensesRepository:           expensesRepository,
	}
}
