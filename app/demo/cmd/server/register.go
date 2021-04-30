package server

import (
	demo_v0 "nautilus/rpc/demo/v0"
	serverDemo_v0 "nautilus/server/demov0"

	"github.com/gin-gonic/gin"
)

func register(router *gin.Engine, internal bool) {

	// 内网接口
	if internal {
		demo_v0.RegisterBlogServiceHTTPServer(router, &serverDemo_v0.DemoServer{})
	}
}
