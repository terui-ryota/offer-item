package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/offer-item/internal/common"
	"github.com/terui-ryota/offer-item/pkg/apperr"
	"github.com/terui-ryota/offer-item/pkg/logger"
	"github.com/terui-ryota/offer-item/servers"
	"github.com/terui-ryota/offer-item/servers/grpc_proxyserver"
	"github.com/terui-ryota/protofiles/go/offer_item"
	"go.opencensus.io/plugin/ocgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func NewApp(
	handler offer_item.OfferItemHandlerServer,
	cfg *config.GRPCConfig,
) common.App {
	opts := []interface{}{
		servers.WithGrpcService(func(s *grpc.Server) {
			// gRPC health check probing endpoint
			grpc_health_v1.RegisterHealthServer(s, health.NewServer())
			offer_item.RegisterOfferItemHandlerServer(s, handler)
		}),
	}
	interceptors := []grpc.UnaryServerInterceptor{
		//common_metadata.UnaryServerInterceptor(),
		apperr.ApplicationErrorUnaryServerInterceptor(),
	}
	opts = append(opts,
		// nantes/libgo で定義されている不要な ocgrpc の server の view が登録されないように無効化
		servers.WithoutOpentracing(),
		servers.WithGrpcRawServerOption(grpc.StatsHandler(&ocgrpc.ServerHandler{})),
		servers.WithAccessLogger(
			[]string{
				servers.GrpcLivenessCheckPath,
				"/grpc.health.v1.Health.Check",
			},
		),
		servers.WithGrpcUnaryInterceptor(interceptors...),
	)

	server, err := servers.NewGrpcServer(cfg.GrpcPort, opts...)
	if err != nil {
		panic(fmt.Errorf("servers.NewGrpcServer: %w", err))
	}
	// grpcサーバのアドレス
	endpoint := fmt.Sprintf("localhost:%d", cfg.GrpcPort)

	// offer-itemはtama-bffでproxyする為、http->grpc変換を行う必要はない。permissionチェックも不要
	proxy, err := grpc_proxyserver.NewGrpcProxyServer(endpoint, nil, nil)
	if err != nil {
		panic(fmt.Errorf("grpc_proxyserver.NewGrpcProxyServer: %w", err))
	}

	return &App{cfg: cfg, server: server, proxy: proxy}
}

type App struct {
	cfg    *config.GRPCConfig
	server servers.Server
	proxy  grpc_proxyserver.GrpcProxyServer
}

//func (a *App) Configure() error {
//	if a.cfg.Tracing != nil {
//		if err := tracing.Configure(a.cfg.Tracing); err != nil {
//			logger.Default().Warn("failed to configure tracing.", zap.Error(err))
//		}
//	}
//	return nil
//}

func (a *App) Start() {
	//defer func() {
	//	_ = tracing.Close()
	//}()
	//if err := a.Configure(); err != nil {
	//	logger.Default().Error("failed to configure app.", zap.Error(err))
	//	panic(err)
	//}

	go func() {
		_ = a.proxy.Start()
	}()

	if err := a.server.Start(context.Background()); err != nil {
		logger.Default().Error("failed to start app.", zap.Error(err))
		panic(err)
	}

	waitForStopSignal()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.proxy.Shutdown(ctx); err != nil {
		logger.Default().Error("failed to stop grpc proxy gracefully.", zap.Error(err))
	}

	if err := a.server.Stop(); err != nil {
		logger.Default().Error("failed to stop app gracefully.", zap.Error(err))
	}
}

func waitForStopSignal() {
	logger.Default().Info("Wait for stop signal")

	sig :=
		waitForSignal(
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM)

	logger.Default().Info("Signal receive", zap.String("signal", sig.String()))
}

func waitForSignal(signals ...os.Signal) os.Signal {
	ch := make(chan os.Signal, 1)

	signal.Notify(ch, signals...)

	return <-ch
}
