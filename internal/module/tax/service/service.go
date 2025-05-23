package service

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *taxService) CreateOrUpdateTax(ctx context.Context, req *dto.CreateOrUpdateTaxRequest) (*dto.CreateOrUpdateTaxResponse, error) {
	tax, err := s.taxRepository.FindTax(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateTax - Failed to find member discount")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateTax - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::CreateOrUpdateTax - Failed to rollback transaction")
			}
		}
	}()

	if tax == nil {
		tax, err = s.taxRepository.InsertNewTax(ctx, tx, &entity.Tax{
			Tax: req.TaxAmount,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::CreateOrUpdateTax - Failed to insert new tax")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	} else {
		tax, err = s.taxRepository.UpdateTax(ctx, tx, &entity.Tax{
			ID:  tax.ID,
			Tax: req.TaxAmount,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::CreateOrUpdateTax - Failed to update member discount")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateTax - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateTaxResponse{
		ID:  tax.ID,
		Tax: fmt.Sprintf("%d", tax.Tax),
	}, nil
}
