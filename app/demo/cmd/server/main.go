package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nautilus/util/conf"
	"nautilus/util/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	reload := make(chan struct{}, 1)
	stop := make(chan os.Signal, 1)

	// 监听配置文件变更
	conf.OnConfigChange(func() { reload <- struct{}{} })
	conf.WatchConfig()
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	go func() {
		for {
			select {
			case <-reload:
				// TODO reset
				os.Exit(0)
			case sg := <-stop:
				fmt.Println("exit ....")
				if sg == syscall.SIGINT {
					os.Exit(0)
				} else {
					os.Exit(0)
				}
			}
		}
	}()

	fmt.Println("start server")
	startServer()
}

func startServer() {
	// TODO ctx 处理
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// middleware
	router.Use(middleware.Timeout(time.Millisecond * 500))
	router.Use(middleware.NewTraceID())

	register(router, internal)
	router.Run("127.0.0.1:8080")
}

func stopServer() {}
