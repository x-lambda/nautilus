package demov0

import (
	"context"
	"fmt"

	"nautilus/util/log"

	pb "nautilus/rpc/demo/v0"
)

type DemoServer struct{}

func (s *DemoServer) CreateArticle(ctx context.Context, req *pb.Article) (resp *pb.Article, err error) {
	// err = demo.TestTimeout(ctx)
	fmt.Println("开始打印日志")
	get := log.Get(ctx)
	if get != nil {
		get.Info("this is q request")
	}
	fmt.Println("打印日志结束")

	// time.Sleep(2 * time.Millisecond)
	resp = &pb.Article{
		Title: "testssss",
	}
	return
}

func (s *DemoServer) GetArticles(ctx context.Context, req *pb.GetArticlesReq) (resp *pb.GetArticlesResp, err error) {
	// demo.TestTimeout(ctx)
	// time.Sleep(100 * time.Millisecond)
	get := log.Get(ctx)
	if get != nil {
		get.Info("this is q request")
	}

	panic("err")

	resp = &pb.GetArticlesResp{
		Total: 10,
	}
	return
}
