package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *ingredientStockHandler) createIngredientStock(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateIngredientStockRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createIngredientStock - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createIngredientStock - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateNewIngredientStock(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createIngredientStock - Failed to create ingredient stock")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *ingredientStockHandler) updateIngredientStock(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateIngredientStockRequest)
		ctx = c.Context()
		id  = c.Params("stock_id")
	)

	stockID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateIngredientStock - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateIngredientStock - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateIngredientStock - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateIngredientStock(ctx, req, stockID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateIngredientStock - Failed to update ingredient stock")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *ingredientStockHandler) getListIngredientStock(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListIngredientStock(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailIngredientStock - Failed to get list ingredient stock")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *ingredientStockHandler) getDetailIngredientStock(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("stock_id")
	)

	stockID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailIngredientStock - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	res, err := h.service.GetDetailIngredientStock(ctx, stockID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailIngredientStock - Failed to get detail ingredient stock")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *ingredientStockHandler) removeIngredientStock(c *fiber.Ctx) error {
	var (
		ctx = c.Context()
		id  = c.Params("stock_id")
	)

	stockID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeIngredientStock - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	err = h.service.RemoveIngredientStock(ctx, stockID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeIngredientStock - Failed to remove ingredient stock")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
