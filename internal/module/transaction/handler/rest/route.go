package rest

import "github.com/gofiber/fiber/v2"

func (h *transactionHandler) TransactionRoute(publicRouter fiber.Router) {
	// public endpoint
	publicRouter.Post("/transactions", h.createTransaction)
}
