package notification

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/notification"
)

type ExternalNotification interface {
	SendTransactionEmail(ctx context.Context, req *notification.SendTransactionEmailRequest) (*notification.SendTransactionEmailResponse, error)
	SendFcmNotification(ctx context.Context, req *notification.SendFcmNotificationRequest) (*notification.SendFcmNotificationResponse, error)
}
