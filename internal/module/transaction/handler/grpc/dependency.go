package grpc

import (
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/cmd/proto/payment"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/adapter"

	// midtransService "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/integration/midtrans/service"
	externalNotification "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/notification"
	externalUserFcmToken "github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/user_fcm_token"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/ports"
	transactionRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/repository"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction/service"
	transactionItemRepository "github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/module/transaction_item/repository"
)

type TransactionAPI struct {
	TransactionService ports.TransactionService
	payment.UnimplementedPaymentCallbackServiceServer
}

func NewTransactionAPI() *TransactionAPI {
	var handler = new(TransactionAPI)

	// external service
	externalNotification := &externalNotification.External{}
	externalUserFcmToken := &externalUserFcmToken.External{}

	// repository
	transactionRepository := transactionRepository.NewTransactionRepository(adapter.Adapters.DzikraPostgres)
	transactionItemRepository := transactionItemRepository.NewTransactionItemRepository(adapter.Adapters.DzikraPostgres)

	// service
	transactionService := service.NewTransactionService(
		adapter.Adapters.DzikraPostgres,
		transactionRepository,
		transactionItemRepository,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		externalNotification,
		externalUserFcmToken,
	)

	// handler
	handler.TransactionService = transactionService

	return handler
}
