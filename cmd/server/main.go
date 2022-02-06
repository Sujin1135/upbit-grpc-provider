package main

import (
	"flag"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	pb "upbit-grpc-provider/cmd/proto"
	"upbit-grpc-provider/internal/config"
)

func main() {
	flag.Parse()
	config.InitLogger()
	config.InitConfig()

	port := 4000
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	} else {
		log.Infof("*** Start to listen the port: %d", port)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterTraderServer(grpcServer, &server{})
	grpcServer.Serve(lis)
}

type server struct {
	pb.UnimplementedTraderServer
}
