package transaction

import (
	"context"
	"fmt"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/constants"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/transaction"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/internal/infrastructure/config"
	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/pkg/utils"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type External struct {
}

func (*External) CreateTransaction(ctx context.Context, req *transaction.CreateTransactionRequest) (*transaction.CreateTransactionResponse, error) {
	conn, err := grpc.Dial(utils.GetEnv("ORDER_GRPC_HOST", config.Envs.Order.OrderGrpcHost), grpc.WithInsecure())
	if err != nil {
		log.Err(err).Msg("external::CreateTransaction - Failed to dial grpc")
		return nil, err
	}
	defer conn.Close()

	client := transaction.NewTransactionServiceClient(conn)

	resp, err := client.CreateTransaction(ctx, req)
	if err != nil {
		log.Err(err).Msg("external::CreateTransaction - Failed to create transaction")
		return nil, err
	}

	if resp.Message != constants.SuccessMessage {
		log.Err(err).Msg("external::CreateTransaction - Response error from order")
		return nil, fmt.Errorf("get response error from order: %s", resp.Message)
	}

	return resp, nil
}
