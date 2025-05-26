package rest

import "github.com/gofiber/fiber/v2"

func (h *transactionHandler) TransactionRoute(publicRouter, superadminRouter fiber.Router) {
	// public endpoint
	publicRouter.Post("/transactions", h.createTransaction)

	// superadmin endpoint
	superadminRouter.Get("/transactions", h.middleware.AuthBearer, h.getListTransaction)
}
