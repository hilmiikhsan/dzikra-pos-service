package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *transactionHandler) createTransaction(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateTransactionRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createTransaction - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createTransaction - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	tableNumber, err := strconv.Atoi(req.TableNumber)
	if err != nil {
		log.Warn().Err(err).Msg("handler::createTransaction - Invalid table number")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid table number"))
	}

	res, err := h.service.CreateTransaction(ctx, req, tableNumber)
	if err != nil {
		log.Error().Err(err).Msg("handler::createTransaction - Failed to create transaction")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}
