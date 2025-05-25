package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	member "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (s *memberDiscountService) CreateOrUpdateMemberDiscount(ctx context.Context, req *dto.CreateOrUpdateMemberDiscountRequest) (*dto.CreateOrUpdateMemberDiscountResponse, error) {
	discount, err := s.memberDiscountRepository.FindDiscount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateMember - Failed to find member discount")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateMember - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::CreateOrUpdateMember - Failed to rollback transaction")
			}
		}
	}()

	if discount == nil {
		discount, err = s.memberDiscountRepository.InsertNewMemberDiscount(ctx, tx, &entity.MemberDiscount{
			Discount: req.Discount,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::CreateOrUpdateMember - Failed to insert new member discount")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	} else {
		discount, err = s.memberDiscountRepository.UpdateMemberDiscount(ctx, tx, &entity.MemberDiscount{
			ID:       discount.ID,
			Discount: req.Discount,
		})
		if err != nil {
			log.Error().Err(err).Msg("service::CreateOrUpdateMember - Failed to update member discount")
			return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
		}
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateOrUpdateMember - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateMemberDiscountResponse{
		ID:        discount.ID,
		Discount:  fmt.Sprintf("%d", discount.Discount),
		UpdatedAt: utils.FormatTime(discount.UpdatedAt),
	}, nil
}

func (s *memberDiscountService) CheckMemberDiscount(ctx context.Context, req *dto.CheckMemberDiscountRequest) (*dto.CheckMemberDiscountResponse, error) {
	memberResult, err := s.memberRepository.FindMemberByEmailOrPhoneNumber(ctx, req.Identifier)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrMemberNotFound) {
			log.Error().Err(err).Msg("service::CheckMemberDiscount - member not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrMemberNotFound))
		}

		log.Error().Err(err).Msg("service::CheckMemberDiscount - Failed to get member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	discount, err := s.memberDiscountRepository.FindFirstMemberDiscount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("service::CheckMemberDiscount - failed to fetch discount")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	resp := &dto.CheckMemberDiscountResponse{
		Discount: fmt.Sprintf("%d", discount.Discount),
		Member: member.GetListMember{
			ID:          memberResult.ID.String(),
			Name:        memberResult.Name,
			Email:       memberResult.Email,
			PhoneNumber: memberResult.PhoneNumber,
			CreatedAt:   utils.FormatTime(memberResult.CreatedAt),
		},
	}

	return resp, nil
}
