package adapter

import (
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

func WithGRPCServer(server *grpc.Server) Option {
	log.Info().Msg("gRPC server connected")
	return func(a *Adapter) {
		a.GRPCServer = server
	}
}
