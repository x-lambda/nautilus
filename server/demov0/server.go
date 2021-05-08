package demov0

import (
	"context"

	"nautilus/service/demo"
	"nautilus/util/log"

	pb "nautilus/rpc/demo/v0"
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
