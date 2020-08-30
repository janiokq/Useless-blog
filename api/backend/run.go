package backend

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/janiokq/Useless-blog/cinit"
	pb "github.com/janiokq/Useless-blog/internal/proto/v1/user"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/janiokq/Useless-blog/internal/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"time"
)

const (
	ServiceName = "api-backend"
)

var (
	UserClient pb.UserServiceClient
)

func Run() {
	//初始化依赖
	cinit.InitOption(ServiceName, cinit.Trace)
	// 建立客户端连接
	grOpts := []grpc_retry.CallOption{
		grpc_retry.WithCodes(codes.Aborted, codes.DeadlineExceeded),
		grpc_retry.WithMax(3),
		grpc_retry.WithPerRetryTimeout(15 * time.Second),
	}

	conn, err := grpc.Dial(cinit.Config.SrvUser.Address,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_opentracing.UnaryClientInterceptor(),
			wrapper.LoggingUnaryClientInterceptor(),
			grpc_retry.UnaryClientInterceptor(grOpts...),
			// wrapper.TraceinUnaryClientInterceptor(),
		)))
	if err != nil {
		logx.Fatal("连接user服务失败" + err.Error())
	}
	// 注册客户端
	UserClient = pb.NewUserServiceClient(conn)

}
