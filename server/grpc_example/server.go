package grpc_example

import (
	"context"

	"nautilus/service/demo"
	"nautilus/util/log"

	pb "nautilus/rpc/grpc/v0"
)

type Server struct {
	// https://github.com/grpc/grpc-go/issues/3794
	pb.UnimplementedGreeterServer
}

// SayHello 实现 SayHello
func (s *Server) SayHello(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloReply, err error) {
	return
}

// SayHelloAgain 实现 SayHelloAgain
func (s *Server) SayHelloAgain(ctx context.Context, req *pb.HelloRequest) (resp *pb.HelloReply, err error) {

	// service 层代码
	err = demo.TestTimeout(ctx)
	log.Get(ctx).Info("this is q request")

	resp = &pb.HelloReply{
		Message: "hihihihihihihi",
	}
	return
}
