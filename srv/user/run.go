package user

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/gateway"
	pb "github.com/janiokq/Useless-blog/internal/proto/v1/user"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/janiokq/Useless-blog/internal/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

const (
	ServiceName = "srv-user"
)

func Run() {
	//初始化需要的依赖

	cinit.InitOption(ServiceName, cinit.Trace, cinit.MySQL, cinit.Redis, cinit.Kafka, cinit.Metrics)
	lis, err := net.Listen("tcp", cinit.Config.SrvUser.Port)
	if err != nil {
		logx.Fatal("failed to listen: " + err.Error())
	}
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_opentracing.UnaryServerInterceptor(),
			wrapper.RecoveryUnaryInterceptor,
			wrapper.LoggingUnaryInterceptor,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_opentracing.StreamServerInterceptor(),
		)),
	)
	pb.RegisterUserServiceServer(s, &Server{})
	reflection.Register(s)
	ctx := context.Background()
	go func() {
		_ = gateway.Run(
			ctx,
			gateway.WithAddr(cinit.Config.SrvUser.GateWayAddr),
			gateway.WithGRPCServer("tcp", cinit.Config.SrvUser.Address),
			gateway.WithSwaggerDir(cinit.Config.SrvUser.GateWaySwaggerDir),
			gateway.WithHandle(pb.RegisterUserServiceHandler),
		)
	}()

	if err := s.Serve(lis); err != nil {
		logx.Fatal("failed to listen: " + err.Error())
	}

	fmt.Println("启动成功")

}
