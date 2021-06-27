package main

import (
	"log"
	"net"

	pb "github.com/husseinelguindi/mtd/protos/mtd"
	"github.com/husseinelguindi/mtd/server"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}

	grpcServer := grpc.NewServer()
	srv := server.NewServer()
	pb.RegisterMtdServer(grpcServer, &srv)

	grpcServer.Serve(lis)
}
