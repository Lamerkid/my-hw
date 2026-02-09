package internalgrpc

import (
	context "context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerRequestLoggerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, _ := metadata.FromIncomingContext(ctx)
		start := time.Now()
		fmt.Printf("%s [%s] %s %s %d %d \"%s\"\n",
			strings.Split(md.Get(":authority")[0], ":")[0],
			start.Format("02/Jan/2006:15:04:05 -0700"),
			info.FullMethod,
			"gRPC",
			http.StatusOK,
			time.Since(start),
			md.Get("user-agent")[0])

		return handler(ctx, req)
	}
}
