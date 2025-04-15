package main

import (
	"context"
	"log"
	"net"
	
	"google.golang.org/grpc"
	pb "go_grpc/grpc-hello-word/proto"
)

const (
	// Port is the port on which the server will listen
	Port = ":50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements the SayHello method of the Greeter service
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}
func main() {
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterGreeterServer(grpcServer, &server{})
	log.Printf("Server is listening on %s", Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}