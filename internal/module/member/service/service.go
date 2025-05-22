package service

import (
	"context"
	"strings"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func (s *memberService) CreateMember(ctx context.Context, req *dto.CreateOrUpdateMemberRequest) (*dto.CreateOrUpdateMemberResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateMember - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::CreateMember - Failed to rollback transaction")
			}
		}
	}()

	id, err := utils.GenerateUUIDv7String()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::CreateMember - Failed to generate UUID")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	payload := &entity.Member{
		ID:          id,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
	}

	// Create new member
	res, err := s.memberRepository.InsertNewMember(ctx, tx, payload)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrEmailAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateMember - Email already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrEmailAlreadyRegistered))
		}

		if strings.Contains(err.Error(), constants.ErrPhoneNumberAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::CreateMember - Phone number already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrPhoneNumberAlreadyRegistered))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::CreateMember - Failed to create new member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::CreateMember - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateMemberResponse{
		ID:          res.ID.String(),
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: res.PhoneNumber,
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *memberService) GetListMember(ctx context.Context, page, limit int, search string) (*dto.GetListMemberResponse, error) {
	// calculate pagination
	currentPage, perPage, offset := utils.Paginate(page, limit)

	// get list member
	members, total, err := s.memberRepository.FindListMember(ctx, perPage, offset, search)
	if err != nil {
		log.Error().Err(err).Msg("service::GetListMember - error getting list member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// check if members is nil
	if members == nil {
		members = []dto.GetListMember{}
	}

	// calculate total pages
	totalPages := utils.CalculateTotalPages(total, perPage)

	// create map response
	response := dto.GetListMemberResponse{
		Members:     members,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
		PageSize:    perPage,
		TotalData:   total,
	}

	// return response
	return &response, nil
}

func (s *memberService) UpdateMember(ctx context.Context, req *dto.CreateOrUpdateMemberRequest, id string) (*dto.CreateOrUpdateMemberResponse, error) {
	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Any("payload", req).Msg("service::UpdateFAQ - Failed to begin transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Any("payload", req).Msg("service::UpdateFAQ - Failed to rollback transaction")
			}
		}
	}()

	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("service::UpdateMember - Failed to parse uuid")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("invalid uuid format"))
	}

	payload := &entity.Member{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Email:       req.Email,
		ID:          uuid,
	}

	// Update member
	res, err := s.memberRepository.UpdateMember(ctx, tx, payload, uuid)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrMemberNotFound) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Member not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrMemberNotFound))
		}

		if strings.Contains(err.Error(), constants.ErrEmailAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Email already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrEmailAlreadyRegistered))
		}

		if strings.Contains(err.Error(), constants.ErrPhoneNumberAlreadyRegistered) {
			log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Phone number already registered")
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(constants.ErrPhoneNumberAlreadyRegistered))
		}

		log.Error().Err(err).Any("payload", req).Msg("service::UpdateMember - Failed to update member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	// commit transaction
	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::UpdateMember - failed to commit transaction")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.CreateOrUpdateMemberResponse{
		ID:          res.ID.String(),
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: res.PhoneNumber,
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *memberService) GetDetailMember(ctx context.Context, id string) (*dto.GetListMember, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("service::GetDetailMember - Failed to parse uuid")
		return nil, err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("invalid uuid format"))
	}

	res, err := s.memberRepository.FindMemberByID(ctx, uuid)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrMemberNotFound) {
			log.Error().Err(err).Msg("service::GetDetailMember - Member not found")
			return nil, err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrMemberNotFound))
		}

		log.Error().Err(err).Msg("service::GetDetailMember - Failed to get member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return &dto.GetListMember{
		ID:          res.ID.String(),
		Name:        res.Name,
		Email:       res.Email,
		PhoneNumber: res.PhoneNumber,
		CreatedAt:   utils.FormatTime(res.CreatedAt),
	}, nil
}

func (s *memberService) RemoveMember(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveMember - Failed to parse uuid")
		return err_msg.NewCustomErrors(fiber.StatusBadRequest, err_msg.WithMessage("invalid uuid format"))
	}

	// Begin transaction
	tx, err := s.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("service::RemoveMember - Failed to begin transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Error().Err(rollbackErr).Msg("service::RemoveMember - Failed to rollback transaction")
			}
		}
	}()

	err = s.memberRepository.SoftDeleteMemberByID(ctx, tx, uuid)
	if err != nil {
		if strings.Contains(err.Error(), constants.ErrMemberNotFound) {
			log.Error().Err(err).Msg("service::RemoveMember - Member not found")
			return err_msg.NewCustomErrors(fiber.StatusNotFound, err_msg.WithMessage(constants.ErrMemberNotFound))
		}

		log.Error().Err(err).Msg("service::RemoveMember - Failed to remove member")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	if err = tx.Commit(); err != nil {
		log.Error().Err(err).Msg("service::RemoveMember - failed to commit transaction")
		return err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return nil
}
