package transaction

import (
	"context"

	"github.com/Digitalkeun-Creative/be-dzikra-pos-service/external/proto/transaction"
)

type ExternalTransaction interface {
	CreateTransaction(ctx context.Context, req *transaction.CreateTransactionRequest) (*transaction.CreateTransactionResponse, error)
}
