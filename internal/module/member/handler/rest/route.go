package rest

import "github.com/gofiber/fiber/v2"

func (h *memberHandler) MemberRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/members", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "member"), h.createMember)
	superadminRouter.Get("/members", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "member"), h.getListMember)
	superadminRouter.Put("/members/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "member"), h.updateMember)
	superadminRouter.Get("/members/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "member"), h.getDetailMember)
	superadminRouter.Delete("/members/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "member"), h.removeMember)
}
