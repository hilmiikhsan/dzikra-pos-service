package route

import (
	ingredientStock "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/ingredient_stock/handler/rest"
	mmember "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member/handler/rest"
	mmemberDiscount "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/member_discount/handler/rest"
	product "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product/handler/rest"
	productCategory "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/product_category/handler/rest"
	recipe "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/recipe/handler/rest"
	tax "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/tax/handler/rest"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func SetupRoutes(app *fiber.App) {
	var (
		// userAPI       = app.Group("/api/users")
		superadminAPI = app.Group("/api/superadmin")
	// publicAPI     = app.Group("/api")
	)

	mmember.NewMemberHandler().MemberRoute(superadminAPI)
	mmemberDiscount.NewMemberDiscountHandler().MemberDiscountRoute(superadminAPI)
	productCategory.NewProductCategoryHandler().ProductCategoryRoute(superadminAPI)
	product.NewProductHandler().ProductRoute(superadminAPI)
	tax.NewTaxHandler().TaxRoute(superadminAPI)
	ingredientStock.NewIngredientStockHandler().IngredientStockRoute(superadminAPI)
	recipe.NewRecipeHandler().RecipeRoute(superadminAPI)

	// fallback route
	app.Use(func(c *fiber.Ctx) error {
		var (
			method = c.Method()                       // get the request method
			path   = c.Path()                         // get the request path
			query  = c.Context().QueryArgs().String() // get all query params
			ua     = c.Get("User-Agent")              // get the request user agent
			ip     = c.IP()                           // get the request IP
		)

		log.Info().
			Str("url", c.OriginalURL()).
			Str("method", method).
			Str("path", path).
			Str("query", query).
			Str("ua", ua).
			Str("ip", ip).
			Msg("Route not found.")
		return c.Status(fiber.StatusNotFound).JSON(response.Error("Route not found"))
	})
}
