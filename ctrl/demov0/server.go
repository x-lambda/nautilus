package demov0

import (
	"context"

	"nautilus/pkg/log"
	"nautilus/svc/demo"

	pb "nautilus/api/demo/v0"
)

type DemoServer struct{}

func (s *DemoServer) CreateArticle(ctx context.Context, req *pb.Article) (resp *pb.Article, err error) {
	err = demo.TestTimeout(ctx)
	log.Get(ctx).Info("this is q request")
	resp = &pb.Article{
		Title: "testssss",
	}
	return
}

func (s *DemoServer) GetArticles(ctx context.Context, req *pb.GetArticlesReq) (resp *pb.GetArticlesResp, err error) {
	err = demo.TestTimeout(ctx)
	log.Get(ctx).Info("this is q request")
	resp = &pb.GetArticlesResp{
		Total: 10,
	}
	return
}
