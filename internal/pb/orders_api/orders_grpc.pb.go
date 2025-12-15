package orders_api

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// RegisterOrdersServiceServer — заглушка. Сейчас не регистрирует сервис в gRPC-сервере.
func RegisterOrdersServiceServer(s *grpc.Server, srv OrdersServiceServer) {
	_ = s
	_ = srv
}

// RegisterOrdersServiceHandlerFromEndpoint — заглушка для регистрации grpc-gateway.
func RegisterOrdersServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
	_ = ctx
	_ = mux
	_ = endpoint
	_ = opts
	return nil
}
