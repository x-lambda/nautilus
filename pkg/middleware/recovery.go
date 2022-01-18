package middleware

import (
	"context"
	"runtime/debug"

	"nautilus/pkg/log"
)

// Recovery panic handler
func Recovery(ctx context.Context) {
	if err := recover(); err != nil {

		stack := string(debug.Stack())
		// 这里为了不影响请求处理流程，使用了协程，原则上每个请求的处理过程中不开goroutine
		log.Get(ctx).Errorf("err: [%+v]\n%s", err, stack)

		// 内部panic返回500错误
	}
}
