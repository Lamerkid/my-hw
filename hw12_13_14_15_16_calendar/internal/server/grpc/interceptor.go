package internalgrpc

import (
	context "context"
	"fmt"
	"net/http"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerRequestLoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	start := time.Now()
	fmt.Printf("%s [%s] %s %s %s %d %d \"%s\"\n",
		md,
		start.Format("02/Jan/2006:15:04:05 -0700"),
		info.FullMethod,
		//
		"gRPC",
		http.StatusOK,
		time.Since(start),
		md.Get("user-agent")[0])

	return handler(ctx, req)
}
