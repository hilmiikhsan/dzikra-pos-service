package utils

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

func HandleInsertUniqueError(err error, data interface{}, uniqueConstraints map[string]string) (interface{}, error) {
	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
		if customMessage, exists := uniqueConstraints[pqErr.Constraint]; exists {
			log.Warn().Err(err).Any("payload", data).Msgf("repository::Insert - %s", customMessage)
			return nil, err_msg.NewCustomErrors(fiber.StatusConflict, err_msg.WithMessage(customMessage))
		}

		log.Error().Err(err).Any("payload", data).Msg("repository::Insert - Unknown unique constraint violation")
		return nil, err_msg.NewCustomErrors(fiber.StatusInternalServerError, err_msg.WithMessage(constants.ErrInternalServerError))
	}

	log.Error().Err(err).Any("payload", data).Msg("repository::Insert - Failed to insert data")
	return nil, err
}
