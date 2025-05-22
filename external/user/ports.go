package user

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	user "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/user/token_validation"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type ExternalTokenValidation interface {
	ValidateToken(ctx context.Context, token string) (*GetTokekenValidationResponse, error)
}

type External struct {
}

func (*External) ValidateToken(ctx context.Context, token string) (*GetTokekenValidationResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::ValidateToken - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := user.NewTokenValidationClient(conn)

	resp, err := client.ValidateToken(ctx, &user.TokenRequest{
		Token: token,
	})
	if err != nil {
		log.Err(err).Msg("external::ValidateToken - Failed to validate token")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::ValidateToken - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	var userRoles []ApplicationPermission
	for _, appPerm := range resp.Data.UserRoles {
		var perms []UserRoleAppPermission
		for _, p := range appPerm.ApplicationPermissions {
			perms = append(perms, UserRoleAppPermission{
				ApplicationID: p.ApplicationId,
				Name:          p.Name,
				Permissions:   p.Permissions,
			})
		}

		userRoles = append(userRoles, ApplicationPermission{
			ApplicationPermission: perms,
			Roles:                 appPerm.Roles,
		})
	}

	return &GetTokekenValidationResponse{
		UserID:   resp.Data.UserId,
		Email:    resp.Data.Email,
		FullName: resp.Data.FullName,
		UserRole: userRoles,
	}, nil
}
