package demov0

import (
	"context"
	"time"

	"nautilus/service/demo"

	pb "nautilus/rpc/demo/v0"
)

type DemoServer struct{}

func (s *DemoServer) CreateArticle(ctx context.Context, req *pb.Article) (resp *pb.Article, err error) {
	err = demo.TestTimeout(ctx)
	time.Sleep(10 * time.Second)
	resp = &pb.Article{
		Title: "testssss",
	}
	return
}

func (s *DemoServer) GetArticles(ctx context.Context, req *pb.GetArticlesReq) (resp *pb.GetArticlesResp, err error) {
	demo.TestTimeout(ctx)
	time.Sleep(10 * time.Second)
	resp = &pb.GetArticlesResp{
		Total: 10,
	}
	return
}
