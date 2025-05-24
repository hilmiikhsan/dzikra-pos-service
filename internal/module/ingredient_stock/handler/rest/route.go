package rest

import "github.com/gofiber/fiber/v2"

func (h *ingredientStockHandler) IngredientStockRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/stock/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "ingredient"), h.createIngredientStock)
	superadminRouter.Patch("/stock/update/:stock_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "ingredient"), h.updateIngredientStock)
	superadminRouter.Get("/stock", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "ingredient"), h.getListIngredientStock)
	superadminRouter.Get("/stock/:stock_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "ingredient"), h.getDetailIngredientStock)
	superadminRouter.Delete("/stock/remove/:stock_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "ingredient"), h.removeIngredientStock)
}
