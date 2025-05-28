package user_fcm_token

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/user_fcm_token"
)

type ExternalUserFcmToken interface {
	GetUserFcmTokenByUserID(ctx context.Context, req *user_fcm_token.GetUserFcmTokenByUserIDRequest) (*user_fcm_token.GetUserFcmTokenByUserIDResponse, error)
}
