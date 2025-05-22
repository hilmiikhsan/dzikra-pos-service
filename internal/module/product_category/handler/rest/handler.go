package rest

import (
	"strconv"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/dto"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *productCategoryHandler) createProductCategory(c *fiber.Ctx) error {
	var (
		req = new(dto.CreateOrUpdateProductCategoryRequest)
		ctx = c.Context()
	)

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::createProductCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.CreateOrUpdateProductCategory(ctx, req)
	if err != nil {
		log.Error().Err(err).Msg("handler::createProductCategory - Failed to create product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) updateProductCategory(c *fiber.Ctx) error {
	var (
		req           = new(dto.CreateOrUpdateProductCategoryRequest)
		ctx           = c.Context()
		categoryIDStr = c.Params("category_id")
	)

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	if err := c.BodyParser(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Failed to parse request body")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Failed to parse request body"))
	}

	if err := h.validator.Validate(req); err != nil {
		log.Warn().Err(err).Msg("handler::updateProductCategory - Invalid request body")
		code, errs := err_msg.Errors(err, req)
		return c.Status(code).JSON(response.Error(errs))
	}

	res, err := h.service.UpdateProductCategory(ctx, req, categoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::updateProductCategory - Failed to update product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) getDetailProductCategory(c *fiber.Ctx) error {
	var (
		ctx           = c.Context()
		categoryIDStr = c.Params("category_id")
	)

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::getDetailProductCategory - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	res, err := h.service.GetDetailProductCategory(ctx, categoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDetailProductCategory - Failed to get detail product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) getListProductCategory(c *fiber.Ctx) error {
	var (
		ctx    = c.Context()
		page   = c.QueryInt("page", 1)
		limit  = c.QueryInt("limit", 10)
		search = c.Query("search", "")
	)

	res, err := h.service.GetListProductCategory(ctx, page, limit, search)
	if err != nil {
		log.Error().Err(err).Msg("handler::getListProductCategory - Failed to get list product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}

func (h *productCategoryHandler) removeProductCategory(c *fiber.Ctx) error {
	var (
		ctx           = c.Context()
		categoryIDStr = c.Params("category_id")
	)

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		log.Warn().Err(err).Msg("handler::removeProductCategory - Invalid id format")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("Invalid id format"))
	}

	err = h.service.RemoveProductCategory(ctx, categoryID)
	if err != nil {
		log.Error().Err(err).Msg("handler::removeProductCategory - Failed to remove product category")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success("OK", ""))
}
