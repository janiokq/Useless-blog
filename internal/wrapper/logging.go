package wrapper

import (
	"context"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	logx.Infof("gRPC method: %s, req: %v", info.FullMethod, req, ctx)
	resp, err := handler(ctx, req)
	logx.Infof("gRPC method: %s, resp: %v", info.FullMethod, resp, ctx)
	return resp, err
}

func LoggingUnaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(parentCtx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		err := invoker(parentCtx, method, req, reply, cc, opts...)
		logx.Infof("method: %s,req: %+v, resp: %+v", method, req, reply, parentCtx)
		return err
	}
}
