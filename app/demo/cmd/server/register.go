package server

import (
	demoV0 "github.com/x-lambda/nautilus/rpc/demo/v0"
	serverDemoV0 "github.com/x-lambda/nautilus/server/demov0"

	"github.com/gin-gonic/gin"
)

func register(router *gin.Engine, internal bool) {

	// 内网接口
	if internal {
		demoV0.RegisterBlogServiceHTTPServer(router, &serverDemoV0.DemoServer{})
	}
}
