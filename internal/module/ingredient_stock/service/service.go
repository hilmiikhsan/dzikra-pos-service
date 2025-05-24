package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *ingredientStockService) CreateNewIngredientStock(ctx context.Context, req *dto.CreateIngredientStockRequest) (*dto.CreateIngredientStockResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateNewIngredientStock - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::CreateNewIngredientStock - Failed to rollback transaction")
			}
		}
	}()

	amountPriceRequiredStock := req.PricePerAmountStock * req.RequiredStock

	payload := &entity.IngredientStock{
		Name:                     req.Name,
		Unit:                     req.Unit,
		PricePerAmountStock:      int(req.PricePerAmountStock),
		RequiredStock:            req.RequiredStock,
		AmountPriceRequiredStock: int(amountPriceRequiredStock),
		AmountStockPerPrice:      int(req.AmountStockPerPrice),
	}

	res, err := s.ingredientStockRepository.InsertNewIngredientStock(ctx, tx, payload)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateNewIngredientStock - Failed to create new ingredient stock")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateNewIngredientStock - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateIngredientStockResponse{
		ID:                  res.ID,
		Name:                res.Name,
		RequiredStock:       res.RequiredStock,
		Unit:                res.Unit,
		PricePerAmountStock: res.PricePerAmountStock,
		AmountStockPerPrice: res.AmountStockPerPrice,
	}, nil
}

func (s *ingredientStockService) UpdateIngredientStock(ctx context.Context, req *dto.CreateIngredientStockRequest, id int) (*dto.CreateIngredientStockResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateIngredientStock - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::UpdateIngredientStock - Failed to rollback transaction")
			}
		}
	}()

	amountPriceRequiredStock := req.PricePerAmountStock * req.RequiredStock

	payload := &entity.IngredientStock{
		ID:                       id,
		Name:                     req.Name,
		Unit:                     req.Unit,
		PricePerAmountStock:      int(req.PricePerAmountStock),
		RequiredStock:            req.RequiredStock,
		AmountPriceRequiredStock: int(amountPriceRequiredStock),
		AmountStockPerPrice:      int(req.AmountStockPerPrice),
	}

	res, err := s.ingredientStockRepository.UpdateIngredientStock(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrIngredientStockNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Ingredient stock not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrIngredientStockNotFound))
		}

		log.Error().Err(err).Msg("service::UpdateIngredientStock - Failed to update ingredient stock")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateIngredientStock - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateIngredientStockResponse{
		ID:                  res.ID,
		Name:                res.Name,
		RequiredStock:       res.RequiredStock,
		Unit:                res.Unit,
		PricePerAmountStock: res.PricePerAmountStock,
		AmountStockPerPrice: res.AmountStockPerPrice,
	}, nil
}

func (s *ingredientStockService) GetListIngredientStock(ctx context.Context, page, limit int, search string) (*dto.GetListIngredientStockResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list ingredient stock
	ingredientStocks, total, err := s.ingredientStockRepository.FindListIngredientStock(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListIngredientStock - error getting list ingredient stock")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if ingredient stocks is nil
	if ingredientStocks == nil {
		ingredientStocks = []dto.GetListIngredientStock{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListIngredientStockResponse{
		Stock:       ingredientStocks,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *ingredientStockService) GetDetailIngredientStock(ctx context.Context, id int) (*dto.GetListIngredientStock, error) {
	res, err := s.ingredientStockRepository.FindIngredientStockByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrIngredientStockNotFound) {
			log.Error().Err(err).Msg("service::GetDetailIngredientStock - Ingredient stock not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrIngredientStockNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailIngredientStock - Failed to get ingredient stock")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetListIngredientStock{
		ID:                  res.ID,
		Name:                res.Name,
		RequiredStock:       res.RequiredStock,
		Unit:                res.Unit,
		PricePerAmountStock: res.PricePerAmountStock,
		AmountStockPerPrice: res.AmountStockPerPrice,
	}, nil
}

func (s *ingredientStockService) RemoveIngredientStock(ctx context.Context, id int) error {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveIngredientStock - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveIngredientStock - Failed to rollback transaction")
			}
		}
	}()

	err = s.ingredientStockRepository.SoftDeleteIngredientStockByID(ctx, tx, id)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrIngredientStockNotFound) {
			log.Error().Err(err).Msg("service::RemoveIngredientStock - Ingredient stock not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrIngredientStockNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveIngredientStock - Failed to remove ingredient stock")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveIngredientStock - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
