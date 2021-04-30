package demov0

import (
	"context"

	"github.com/x-lambda/nautilus/service/demo"

	pb "github.com/x-lambda/nautilus/rpc/demo/v0"
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
