package main

import (
	"net"

	pb "nautilus/rpc/grpc/v0"
	"nautilus/server/grpc_example"

	"google.golang.org/grpc"
)

func main() {
	grpcServer := grpc.NewServer()

	// 注册服务
	pb.RegisterGreeterServer(grpcServer, new(grpc_example.Server))

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}

	grpcServer.Serve(lis)
}
