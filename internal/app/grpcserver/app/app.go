package app

import (
	"fmt"
	"log"
	"net"

	"github.com/terui-ryota/offer-item/internal/app/grpcserver/config"
	"github.com/terui-ryota/protofiles/go/offer_item"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/terui-ryota/offer-item/internal/common"
)

func NewApp(
	handler offer_item.OfferItemHandlerServer,
	cfg *config.GRPCConfig,
) common.App {
	return &App{
		cfg:     cfg,
		handler: handler,
	}
}

type App struct {
	cfg     *config.GRPCConfig
	handler offer_item.OfferItemHandlerServer
	server  *grpc.Server
}

func (a *App) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	a.server = grpc.NewServer()

	// gRPC health check probing endpoint
	grpc_health_v1.RegisterHealthServer(a.server, health.NewServer())

	// Register OfferItemHandler
	offer_item.RegisterOfferItemHandlerServer(a.server, a.handler)

	log.Printf("gRPC server listening on port %d", a.cfg.GrpcPort)

	if err := a.server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (a *App) Stop() {
	log.Println("Stopping gRPC server")

	if a.server != nil {
		a.server.GracefulStop()
	}
}
