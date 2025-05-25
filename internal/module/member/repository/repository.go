package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/entity"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/ports"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

var _ ports.MemberRepository = &memberRepository{}

type memberRepository struct {
	db *sqlx.DB
}

func NewMemberRepository(db *sqlx.DB) *memberRepository {
	return &memberRepository{
		db: db,
	}
}

func (r *memberRepository) InsertNewMember(ctx context.Context, tx *sqlx.Tx, data *entity.Member) (*entity.Member, error) {
	var res = new(entity.Member)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryInsertNewMember),
		data.ID,
		data.Name,
		data.PhoneNumber,
		data.Email,
	).Scan(
		&res.ID,
		&res.Name,
		&res.PhoneNumber,
		&res.Email,
		&res.CreatedAt,
	)
	if err != nil {
		uniqueConstraints := map[string]string{
			"members_email_key":        constants.ErrEmailAlreadyRegistered,
			"members_phone_number_key": constants.ErrPhoneNumberAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::InsertNewMember - Failed to insert new member")
			return nil, handleErr
		}

		if member, ok := val.(*entity.Member); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewMember - Failed to insert new member")
			return member, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::InsertNewMember - Failed to insert new member")
		return nil, err
	}

	return res, nil
}

func (r *memberRepository) FindListMember(ctx context.Context, limit, offset int, search string) ([]dto.GetListMember, int, error) {
	args := []interface{}{search, search, search, limit, offset}

	var responses []entity.Member
	if err := r.db.SelectContext(ctx, &responses, r.db.Rebind(queryFindListMember), args...); err != nil {
		log.Error().Err(err).Msg("repository::FindListMember - error executing query")
		return nil, 0, err
	}

	var total int
	if err := r.db.GetContext(ctx, &total, r.db.Rebind(queryCountFindListMember), args[:3]...); err != nil {
		log.Error().Err(err).Msg("repository::FindListMember - error counting members")
		return nil, 0, err
	}

	members := make([]dto.GetListMember, 0, len(responses))
	for _, v := range responses {
		members = append(members, dto.GetListMember{
			ID:          v.ID.String(),
			Name:        v.Name,
			Email:       v.Email,
			PhoneNumber: v.PhoneNumber,
			CreatedAt:   utils.FormatTime(v.CreatedAt),
		})
	}

	return members, total, nil
}

func (r *memberRepository) UpdateMember(ctx context.Context, tx *sqlx.Tx, data *entity.Member, id uuid.UUID) (*entity.Member, error) {
	var res = new(entity.Member)

	err := tx.QueryRowContext(ctx, r.db.Rebind(queryUpdateMember),
		data.Name,
		data.PhoneNumber,
		data.Email,
		data.ID,
	).Scan(
		&res.ID,
		&res.Name,
		&res.PhoneNumber,
		&res.Email,
		&res.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			errMessage := fmt.Errorf("repository::UpdateMember - member with id %d is not found", data.ID)
			log.Error().Err(err).Msg(errMessage.Error())
			return nil, errors.New(constants.ErrMemberNotFound)
		}

		uniqueConstraints := map[string]string{
			"members_email_key":        constants.ErrEmailAlreadyRegistered,
			"members_phone_number_key": constants.ErrPhoneNumberAlreadyRegistered,
		}

		val, handleErr := utils.HandleInsertUniqueError(err, data, uniqueConstraints)
		if handleErr != nil {
			log.Error().Err(handleErr).Any("payload", data).Msg("repository::UpdateMember - Failed to update new member")
			return nil, handleErr
		}

		if member, ok := val.(*entity.Member); ok {
			log.Error().Err(err).Any("payload", data).Msg("repository::UpdateMember - Failed to update new member")
			return member, nil
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::UpdateMember - Failed to update member")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	return res, nil
}

func (r *memberRepository) FindMemberByID(ctx context.Context, id uuid.UUID) (*entity.Member, error) {
	var res = new(entity.Member)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindMemberByID), id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindMemberByID - member not found")
			return nil, errors.New(constants.ErrMemberNotFound)
		}

		log.Error().Err(err).Msg("repository::FindMemberByID - error executing query")
		return nil, err
	}

	return res, nil
}

func (r *memberRepository) SoftDeleteMemberByID(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	result, err := tx.ExecContext(ctx, tx.Rebind(querySoftDeleteMemberByID), id)
	if err != nil {
		log.Error().Err(err).Any("id", id).Msg("repository::SoftDeleteMemberByID - Failed to soft delete member")
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Error().Err(err).Msg("repository::SoftDeleteMemberByID - Failed to fetch rows affected")
		return err
	}

	if rowsAffected == 0 {
		errNotFound := errors.New(constants.ErrMemberNotFound)
		log.Error().Err(errNotFound).Any("id", id).Msg("repository::SoftDeleteMemberByID - Member not found")
		return errNotFound
	}

	return nil
}

func (r *memberRepository) FindMemberByEmailOrPhoneNumber(ctx context.Context, identifier string) (*entity.Member, error) {
	var res = new(entity.Member)

	err := r.db.GetContext(ctx, res, r.db.Rebind(queryFindMemberByEmailOrPhoneNumber), identifier)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Error().Err(err).Msg("repository::FindMemberByEmailOrPhoneNumber - member not found")
			return nil, errors.New(constants.ErrMemberNotFound)
		}

		log.Error().Err(err).Msg("repository::FindMemberByEmailOrPhoneNumber - error executing query")
		return nil, err
	}

	return res, nil
}
