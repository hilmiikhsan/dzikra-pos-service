package rest

import "github.com/gofiber/fiber/v2"

func (h *productCategoryHandler) ProductCategoryRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/category/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "product_category_pos"), h.createProductCategory)
	superadminRouter.Patch("/category/update/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "product_category_pos"), h.updateProductCategory)
	superadminRouter.Get("/category/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_category_pos"), h.getDetailProductCategory)
	superadminRouter.Get("/category", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_category_pos"), h.getListProductCategory)
	superadminRouter.Delete("/category/remove/:category_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "product_category_pos"), h.removeProductCategory)
}
