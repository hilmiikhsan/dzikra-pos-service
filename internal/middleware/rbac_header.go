package middleware

import (
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *AuthMiddleware) RBACMiddleware(resource, action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		raw, ok := c.Locals("user_roles").([]UserRoleDetail)
		if !ok {
			log.Error().Msg("middleware::RBACMiddleware - user_roles not found or wrong type in context")
			return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
		}

		required := fmt.Sprintf("%s|%s", action, resource)

		has := false
		for _, ur := range raw {
			for _, app := range ur.ApplicationPermission {
				for _, perm := range app.Permissions {
					if perm == required {
						has = true
						break
					}
				}
				if has {
					break
				}
			}
			if has {
				break
			}
		}

		if !has {
			log.Error().
				Str("resource", resource).
				Str("action", action).
				Msg("middleware::RBACMiddleware - User does not have permission")
			return c.Status(fiber.StatusForbidden).JSON(response.Error(constants.ErrApplicationForbidden))
		}

		return c.Next()
	}
}
