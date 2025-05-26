package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/expenses/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *expensesService) CreateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest) (*dto.CreateOrUpdateExpensesResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::CreateExpenses - Failed to rollback transaction")
			}
		}
	}()

	date, err := utils.ParseTime(req.Date)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - Failed to parse date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidDateFormat))
	}

	payload := &entity.Expenses{
		Name: req.Name,
		Cost: req.Cost,
		Date: date,
	}

	res, err := s.expensesRepository.InsertNewExpenses(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - Failed to insert new expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateExpenses - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateExpensesResponse{
		ID:   res.ID,
		Name: res.Name,
		Cost: res.Cost,
		Date: utils.FormatTime(res.Date),
	}, nil
}

func (s *expensesService) GetListExpenses(ctx context.Context, page, limit int, search string) (*dto.GetListExpensesResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list expenses
	expenses, total, err := s.expensesRepository.FindListExpenses(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListMember - error getting list expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if expenses is nil
	if expenses == nil {
		expenses = []dto.GetListExpenses{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListExpensesResponse{
		Expenses:    expenses,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *expensesService) GetDetailExpenses(ctx context.Context, id int) (*dto.GetListExpenses, error) {
	expenses, err := s.expensesRepository.FindExpensesByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrExpensesNotFound) {
			log.Error().Err(err).Msg("service::GetDetailExpenses - Expenses not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrExpensesNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailExpenses - error getting expenses by id")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetListExpenses{
		ID:        expenses.ID,
		Name:      expenses.Name,
		Cost:      expenses.Cost,
		Date:      utils.FormatTime(expenses.Date),
		CreatedAt: utils.FormatTime(expenses.CreatedAt),
		UpdatedAt: utils.FormatTime(expenses.UpdatedAt),
	}, nil
}

func (s *expensesService) UpdateExpenses(ctx context.Context, req *dto.CreateOrUpdateExpensesRequest, id int) (*dto.CreateOrUpdateExpensesResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateExpenses - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::UpdateExpenses - Failed to rollback transaction")
			}
		}
	}()

	date, err := utils.ParseTime(req.Date)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateExpenses - Failed to parse date")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage(constants.ErrInvalidDateFormat))
	}

	payload := &entity.Expenses{
		ID:   id,
		Name: req.Name,
		Cost: req.Cost,
		Date: date,
	}

	res, err := s.expensesRepository.UpdateExpenses(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrExpensesNotFound) {
			log.Error().Err(err).Msg("service::UpdateExpenses - Expenses not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrExpensesNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateExpenses - Failed to update expenses")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateExpenses - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateExpensesResponse{
		ID:   res.ID,
		Name: res.Name,
		Cost: res.Cost,
		Date: utils.FormatTime(res.Date),
	}, nil
}

func (s *expensesService) RemoveExpenses(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveExpenses - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveExpenses - Failed to rollback transaction")
			}
		}
	}()

	err = s.expensesRepository.SoftDeleteExpensesByID(ctx, tx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrExpensesNotFound) {
			log.Error().Err(err).Msg("service::RemoveExpenses - Expenses not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrExpensesNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveExpenses - Failed to remove expenses")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveExpenses - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
