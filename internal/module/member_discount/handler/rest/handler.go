package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *memberDiscountHandler) createOrUpdateMemberDiscount(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateMemberDiscountRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createOrUpdateMemberDiscount - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createOrUpdateMemberDiscount - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateOrUpdateMemberDiscount(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createOrUpdateMemberDiscount - Failed to create or update member discount")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
