package middleware

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *AuthMiddleware) AuthBearer(c *fiber.Ctx) error {
	accessToken := c.Get(constants.HeaderAuthorization)

	// If the cookie is not set, return an unauthorized status
	if accessToken == "" {
		log.Error().Msg("middleware::AuthBearer - Unauthorized [Header not set]")
		c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrAccessTokenIsRequired))
	}

	// remove the Bearer prefix
	if len(accessToken) > 7 {
		accessToken = accessToken[7:]
	}

	claims, err := m.externalAuth.ValidateToken(c.Context(), accessToken)
	if err != nil {
		log.Error().Err(err).Any("payload", accessToken).Msg("middleware::AuthBearer - Error while parsing token")
		return c.Status(fiber.StatusUnauthorized).JSON(response.Error(constants.ErrTokenAlreadyExpired))
	}

	var ur []UserRoleDetail
	for _, cr := range claims.UserRole {
		var apps []ApplicationPermissionDetail
		for _, ap := range cr.ApplicationPermission {
			apps = append(apps, ApplicationPermissionDetail{
				ApplicationID: ap.ApplicationID,
				Name:          ap.Name,
				Permissions:   ap.Permissions,
			})
		}
		ur = append(ur, UserRoleDetail{
			Roles:                 cr.Roles,
			ApplicationPermission: apps,
		})
	}

	c.Locals("user_id", claims.UserID)
	c.Locals("email", claims.Email)
	c.Locals("full_name", claims.FullName)
	c.Locals("user_roles", ur)

	// If the token is valid, pass the request to the next handler
	return c.Next()
}
