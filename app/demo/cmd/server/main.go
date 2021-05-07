package server

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"nautilus/util/conf"
	"nautilus/util/log"
	"nautilus/util/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
				log.Reset()
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

	startServer()
}

func startServer() {
	// gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	// metrics
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// middleware
	router.Use(middleware.Logging())
	router.Use(middleware.Timeout(time.Millisecond * 50000))
	router.Use(middleware.NewTraceID())

	register(router, internal)
	router.Run(":8080")
}

func stopServer() {}
