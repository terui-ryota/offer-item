package main

import (
	"context"
	"log"
	"net"

	pb "github.com/terui-ryota/offer-item/protofiles/proto"

	"google.golang.org/grpc"
)

const (
	port = ":50052"
)

type server struct {
	pb.UnimplementedPrivateServer
}

func (s *server) HealthCheck(ctx context.Context, req *pb.HealthCheckReq) (*pb.HealthCheckRes, error) {
	log.Println(req.CheckStr)
	return &pb.HealthCheckRes{
		Request: req,
		Num:     1,
	}, nil
}

func (s *server) mustEmbedUnimplementedPrivateServer() {}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPrivateServer(s, &server{})
	log.Printf("Server listening on port %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
