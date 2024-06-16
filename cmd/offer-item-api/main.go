package main

import (
	"context"
	"log"
	"net"

	pb "github.com/terui-ryota/protofiles/go/offer_item"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type server struct {
	pb.UnimplementedOfferItemHandlerServer
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	log.Println(req.CheckStr)
	return &pb.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

func (s *server) SaveOfferItem(ctx context.Context, req *pb.SaveOfferItemRequest) (*pb.SaveOfferItemResponse, error) {
	// TODO: SaveOfferItem の実装を追加
	return &pb.SaveOfferItemResponse{}, nil
}

func (s *server) mustEmbedUnimplementedOfferItemHandlerServer() {}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterOfferItemHandlerServer(s, &server{})
	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
