package gateway

import (
	"context"
	gwruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/janiokq/Useless-blog/internal/utils/logx"
	"net/http"
)

const (
	DefaultAddr        = ":8888" //  默认地址
	DefaultGrpcAddr    = ":5001" //  默认grpc地址
	DefaultGrpcNetwork = "tcp"   //  默认grpc
	DefaultSwaggerDir  = "/swagger"
)

type Endpoint struct {
	Network, Addr string
}

type Options struct {
	//  Addr is the address to listen
	Addr string

	//  GRPCServer defines an endpoint of a gRPC service
	GRPCServer Endpoint

	//  SwaggerDir is a path to a directory from which the server
	//  serves swagger specs.
	SwaggerDir string

	//  Mux is a list of options to be passed to the grpc-gateway multiplexer
	Mux []gwruntime.ServeMuxOption

	//  注册
	Handles []regHandle
}

type Option func(*Options)

//  设置监听地址
func WithAddr(addr string) Option {
	return func(opts *Options) {
		opts.Addr = addr
	}
}

//  设置grpc服务地址
func WithGRPCServer(network, addr string) Option {
	return func(opts *Options) {
		opts.GRPCServer = Endpoint{
			Addr:    addr,
			Network: network,
		}
	}
}

//  设置swagger目录
func WithSwaggerDir(dir string) Option {
	return func(opts *Options) {
		opts.SwaggerDir = dir
	}
}

//  设置mux
func WithMuxOption(mux ...gwruntime.ServeMuxOption) Option {
	return func(opts *Options) {
		opts.Mux = mux
	}
}

//  设置Handles
func WithHandle(handle ...regHandle) Option {
	return func(opts *Options) {
		opts.Handles = handle
	}
}

//  开启一个网关服务
/* go gateway.Run(
ctx,
gateway.WithAddr(":8888"),
gateway.WithGRPCServer("tcp", ":5001"),
gateway.WithSwaggerDir("/swagger"),
gateway.WithHandle(pb.RegisterRbacServiceHandler),
)
*/
//  Run starts a HTTP server and blocks while running if successful.//  The server will be shutdown when "ctx" is canceled.
func Run(ctx context.Context, options ...Option) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//  初始化参数
	opts := &Options{
		Addr: DefaultAddr,
		GRPCServer: Endpoint{
			Network: DefaultGrpcNetwork,
			Addr:    DefaultGrpcAddr,
		},
		SwaggerDir: DefaultSwaggerDir,
		Mux:        make([]gwruntime.ServeMuxOption, 0),
		Handles:    make([]regHandle, 0),
	}
	for _, f := range options {
		f(opts)
	}
	conn, err := dial(ctx, opts.GRPCServer.Network, opts.GRPCServer.Addr)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		if err = conn.Close(); err != nil {
			logx.Errorf("Failed to close a client connection to the gRPC server: %v", err)
		}
	}()

	mux := http.NewServeMux()
	mux.HandleFunc("/swagger/", swaggerServer(opts.SwaggerDir))
	mux.HandleFunc("/healthz", healthzServer(conn))

	gw, err := NewGateway(ctx, conn, opts.Mux, opts.Handles)
	if err != nil {
		return err
	}
	mux.Handle("/", gw)

	s := &http.Server{
		Addr:    opts.Addr,
		Handler: allowCORS(mux),
	}
	go func() {
		<-ctx.Done()
		logx.Infof("Shutting down the http server")
		if err := s.Shutdown(context.Background()); err != nil {
			logx.Errorf("Failed to shutdown http server: %v", err)
		}
	}()

	logx.Infof("Starting listening at %s", opts.Addr)
	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		logx.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil
}
