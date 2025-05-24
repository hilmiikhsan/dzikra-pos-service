package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *recipeHandler) getListRecipe(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListRecipe(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListRecipe - Failed to get list recipe")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *recipeHandler) updateRecipe(c *fiber.Ctx) error {
	var (
		req = new(dto.UpdateRecipeRequest)
		ctx = c.Context()
		id  = c.Params("product_id")
	)

	productID, err := strconv.Atoi(id)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateRecipe - Invalid product_id")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid product_id"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateRecipe - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateRecipe - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateRecipe(ctx, req, productID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateRecipe - Failed to update recipe")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
