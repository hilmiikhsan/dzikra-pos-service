package rest

import "github.com/gofiber/fiber/v2"

func (h *dashboardHandler) DashboardRoute(superadminRouter fiber.Router) {
	superadminRouter.Get("/dashboard", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "dashboard_pos"), h.getDashboard)
}
