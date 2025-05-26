package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func (h *transactionHandler) getListTransaction(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListTransaction(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListTransaction - Failed to get list of transactions")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *transactionHandler) getTransactionDetail(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("id")
	)

	if id == ":id" {
		log.Warn().Msg("handler::getTransactionDetail - Transaction ID is required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Transaction ID is required"))
	}

	_, err := uuid.Parse(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getTransactionDetail - Invalid transaction ID format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid transaction ID format"))
	}

	res, err := h.service.GetTransactionDetail(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("handler::getTransactionDetail - Failed to get transaction detail")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
