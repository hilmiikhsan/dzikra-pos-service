package rest

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/err_msg"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *dashboardHandler) getDashboard(c *fiber.Ctx) error {
	var (
		ctx       = c.Context()
		startDate = c.Query("start_date")
		endDate   = c.Query("end_date")
	)

	if startDate == "" || endDate == "" {
		log.Warn().Msg("handler::getDashboard - start_date and end_date query params are required")
		return c.Status(fiber.StatusBadRequest).JSON(response.Error("start_date and end_date query params are required"))
	}

	res, err := h.service.GetDashboard(ctx, startDate, endDate)
	if err != nil {
		log.Error().Err(err).Msg("handler::getDashboard - Failed to get dashboard")
		code, errs := err_msg.Errors[error](err)
		return c.Status(code).JSON(response.Error(errs))
	}

	return c.Status(fiber.StatusOK).JSON(response.Success(res, ""))
}
