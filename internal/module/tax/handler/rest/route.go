package rest

import "github.com/gofiber/fiber/v2"

func (h *taxHandler) TaxRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Patch("/tax", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "tax"), h.createOrUpdateTax)
}
