package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	pb "upbit-grpc-provider/cmd/proto"
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 3000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTraderServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

type server struct {
	pb.UnimplementedTraderServer
}
