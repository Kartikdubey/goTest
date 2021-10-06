package main

import (
	"log"
	"net"

	"github.com/Kartikdubey/goTest/clientserver"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := clientserver.Server{}
	grpcServer := grpc.NewServer()

	clientserver.RegisterServiceOneServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
