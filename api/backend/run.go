package backend

import (
	"github.com/gin-gonic/gin"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/janiokq/Useless-blog/cinit"
	"github.com/janiokq/Useless-blog/internal/api"
	"github.com/janiokq/Useless-blog/internal/metrics"
	pb "github.com/janiokq/Useless-blog/internal/proto/v1/user"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"github.com/janiokq/Useless-blog/internal/wrapper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net/http"
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
	cinit.InitOption(ServiceName, cinit.Trace, cinit.Redis)
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
	outer := gin.New()
	outer.Use(gin.Recovery())
	outer.Use(gin.Logger())
	//设置允许跨域
	outer.Use(api.Cors())
	outer.Use(api.TraceHeader())
	// Metrics
	if cinit.Config.Metrics.Enable == "yes" {
		//  Push模式
		m := metrics.NewMetrics()
		outer.Use(api.MetricsFunc(m))
		m.MemStats()
		//  InfluxDB
		m.InfluxDBWithTags(
			time.Duration(cinit.Config.Metrics.Duration)*time.Second,
			cinit.Config.Metrics.URL,
			cinit.Config.Metrics.Database,
			cinit.Config.Metrics.UserName,
			cinit.Config.Metrics.Password,
			cinit.Config.Metrics.Measurement,
			map[string]string{"service": ServiceName},
		)
	}

	// 总分组 需要 token 接口
	g := outer.Group("/backend/v1", api.JWT())

	// 用户
	g.GET("/user", UserInfo)
	g.PUT("/user", UserUpdate)

	// 总分组  公开接口
	gu := outer.Group("/frontend/v1")
	gu.POST("/login", Login)
	gu.POST("/register", Register)

	// check
	check := outer.Group("/backend/check")
	check.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	// 启动service
	if err := outer.Run(cinit.Config.APIBackend.Port); err != nil {
		logx.Fatal("启动服务失败" + err.Error())
	}
	defer conn.Close()

}
