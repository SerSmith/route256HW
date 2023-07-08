package callCounter

import (
	"context"
	"google.golang.org/grpc"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	h, err := handler(ctx, req)

	HandlerCounter.WithLabelValues(info.FullMethod).Inc()

	return h, err
}
