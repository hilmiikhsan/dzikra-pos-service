package service

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/dashboard/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *dashboardService) GetDashboard(ctx context.Context, startDate, endDate string) (*dto.GetDashboardResponse, error) {
	start, err := utils.ParseDateToUTC(startDate)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error parsing start date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidStartDate))
	}

	end, err := utils.ParseEndDateToUTC(endDate)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error parsing end date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidEndDate))
	}

	totalAmount, err := s.transactionHistoryRepository.SumTotalAmount(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total amount transaction history")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalCostExpenses, err := s.expensesRepository.SumTotalCostExpenses(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total cost expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalCountTransactionHistory, err := s.transactionHistoryRepository.CountTransactionHistory(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total count transaction history")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalQuantity, err := s.transactionHistoryRepository.SumTotalQuantity(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total quantity transaction history")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	totalCapital, err := s.transactionHistoryRepository.SumTotalCapital(ctx, start, end)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDashboard - error getting total capital transaction history")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	netSales := totalAmount - totalCapital
	profitLoss := netSales - totalCostExpenses

	return &dto.GetDashboardResponse{
		TotalAmount:       totalAmount,
		TotalExpenses:     totalCostExpenses,
		TotalTransactions: totalCountTransactionHistory,
		TotalSelling:      totalQuantity,
		TotalCapital:      totalCapital,
		NetSales:          netSales,
		ProfitLoss:        profitLoss,
	}, nil
}
