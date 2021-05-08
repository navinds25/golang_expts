package main

import (
	"net"

	"golang.org/x/net/context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "github.com/navinds25/grpcGoExpts/eaconnproto"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello" + in.Name}, nil
}

func (s *server) SendConfig(ctx context.Context, in *pb.Config) (*pb.Ack, error) {
	return &pb.Ack{Valid: "true"}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	pb.RegisterPublisherServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
