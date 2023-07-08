package timeHist

import (
	"context"
	"time"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc"
)

func MiddlewareGRPC(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	
	start := time.Now()

	h, err := handler(ctx, req)

    elapsed := time.Since(start)


	histogramByGroup.WithLabelValues(status.Code(err).String(), info.FullMethod).Observe(elapsed.Seconds())

	return h, err
}
