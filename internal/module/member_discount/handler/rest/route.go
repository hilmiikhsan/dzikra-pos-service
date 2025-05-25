package rest

import "github.com/gofiber/fiber/v2"

func (h *memberDiscountHandler) MemberDiscountRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Put("/member-discount", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create_update", "member_discount"), h.createOrUpdateMemberDiscount)
	superadminRouter.Post("/member-discount/check", h.checkMemberDiscount)
}
