package user_fcm_token

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/user_fcm_token"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) GetUserFcmTokenByUserID(ctx context.Context, req *user_fcm_token.GetUserFcmTokenByUserIDRequest) (*user_fcm_token.GetUserFcmTokenByUserIDResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("AUTH_GRPC_HOST", config.Envs.Auth.AuthGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::GetUserFcmTokenByUserID - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := user_fcm_token.NewUserFcmTokenClient(conn)

	resp, err := client.GetUserFcmTokenByUserID(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::GetUserFcmTokenByUserID - Failed to get user fcm token by user id")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::GetUserFcmTokenByUserID - Response error from auth")
		return nil, fmt.Errorf("get response error from auth: %s", resp.Message)
	}

	return resp, nil
}
