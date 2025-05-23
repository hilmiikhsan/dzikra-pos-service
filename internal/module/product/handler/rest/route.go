package rest

import "github.com/gofiber/fiber/v2"

func (h *productHandler) ProductRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/product/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "product_pos"), h.createProduct)
	superadminRouter.Patch("/product/update/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "product_pos"), h.updateProduct)
	superadminRouter.Get("/product", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_pos"), h.getListProduct)
	superadminRouter.Get("/product/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "product_pos"), h.getDetailProduct)
	superadminRouter.Delete("/product/remove/:product_id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "product_pos"), h.removeProduct)
}
