package rest

import "github.com/gofiber/fiber/v2"

func (h *expensesHandler) ExpensesRoute(superadminRouter fiber.Router) {
	// superadmin endpoint
	superadminRouter.Post("/expenses/create", h.middleware.AuthBearer, h.middleware.RBACMiddleware("create", "expenses"), h.createExpenses)
	superadminRouter.Get("/expenses", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "expenses"), h.getListExpenses)
	superadminRouter.Get("/expenses/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("read", "expenses"), h.getDetailExpenses)
	superadminRouter.Patch("/expenses/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("update", "expenses"), h.updateExpenses)
	superadminRouter.Delete("/expenses/:id", h.middleware.AuthBearer, h.middleware.RBACMiddleware("delete", "expenses"), h.removeExpenses)
}
