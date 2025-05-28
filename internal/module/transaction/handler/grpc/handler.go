package grpc

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/cmd/proto/payment"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/dto"
	"github.com/rs/zerolog/log"
)

func (api *TransactionAPI) CallbackPayment(ctx context.Context, req *payment.PaymentCallbackRequest) (*payment.PaymentCallbackResponse, error) {
	err := api.TransactionService.CallbackPayment(ctx, &dto.PaymentCallbackRequest{
		PaymentID:     req.PaymentId,
		TransactionID: req.TransactionId,
		Status:        req.Status,
		UserFcmToken:  req.UserFcmToken,
		UserID:        req.UserId,
		FullName:      req.FullName,
		Email:         req.Email,
	})
	if err != nil {
		log.Err(err).Msg("order::CallbackPayment - Failed to callback payment")
		return &payment.PaymentCallbackResponse{
			Message: "failed to create order",
		}, nil
	}

	return &payment.PaymentCallbackResponse{
		Message: "success",
	}, nil
}
