package demov0

import (
	"context"

	pb "nautilus/rpc/demo/v0"
	"nautilus/service/demo"
)

type DemoServer struct{}

func (s *DemoServer) CreateArticle(ctx context.Context, req *pb.Article) (resp *pb.Article, err error) {
	err = demo.TestTimeout(ctx)
	resp = &pb.Article{}
	return
}

func (s *DemoServer) GetArticles(ctx context.Context, req *pb.GetArticlesReq) (resp *pb.GetArticlesResp, err error) {
	demo.TestTimeout(ctx)
	resp = &pb.GetArticlesResp{}
	return
}
