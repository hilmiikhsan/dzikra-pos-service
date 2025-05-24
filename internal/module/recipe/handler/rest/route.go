package rest

import "github.com/gofiber/fiber/v2"

func (h *recipeHandler) RecipeRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Get("/recipe", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "recipe"), h.getListRecipe)
	superadminRouter.Patch("/recipe/update/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "recipe"), h.updateRecipe)
}
