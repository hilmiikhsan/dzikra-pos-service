package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *memberHandler) createMember(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateMemberRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createMember - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createMember - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateMember(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createMember - Failed to create member")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *memberHandler) getListMember(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListMember(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListMember - Failed to get list member")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *memberHandler) updateMember(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateMemberRequest)
		ctx = c.Context()
		id  = c.Params("id")
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateMember - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateMember - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateMember(ctx, req, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateMember - Failed to update member")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *memberHandler) getDetailMember(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("id")
	)

	res, err := h.service.GetDetailMember(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailMember - Failed to get detail member")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *memberHandler) removeMember(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("id")
	)

	err := h.service.RemoveMember(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeMember - Failed to remove member")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
