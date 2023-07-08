package logger

import (
	"context"

	"google.golang.org/grpc"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	h, err := handler(ctx, req)
	if err != nil {
		Errorf(ctx, info.FullMethod, "error while processing handler, err=%s", err)
	}

	return h, err
}

