package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// // applyReadTimeout 应用读超时逻辑
// func applyReadTimeout(c *gin.Context, readTimeout time.Duration) bool {
// 	readCtx, readCancel := context.WithTimeout(c.Request.Context(), readTimeout)
// 	defer readCancel()
// 	c.Request = c.Request.WithContext(readCtx)

// 	done := make(chan struct{})
// 	go func() {
// 		c.Next()
// 		close(done)
// 	}()

// 	select {
// 	case <-done:
// 		return false
// 	case <-readCtx.Done():
// 		c.AbortWithStatus(http.StatusRequestTimeout)
// 		return true
// 	}
// }

// // applyWriteTimeout 应用写超时逻辑
// func applyWriteTimeout(c *gin.Context, writeTimeout time.Duration) {
// 	writeCtx, writeCancel := context.WithTimeout(context.Background(), writeTimeout)
// 	defer writeCancel()

// 	done := make(chan struct{})
// 	go func() {
// 		// 模拟写操作完成
// 		close(done)
// 	}()

// 	select {
// 	case <-done:
// 		// 正常完成
// 	case <-writeCtx.Done():
// 		// 写超时
// 		c.Writer.WriteHeader(http.StatusRequestTimeout)
// 	}
// }

// TimeoutMiddleware 返回一个Gin中间件，用于处理连接的超时空值
// 对sse连接不设置超时，对其他连接设置超时时间
func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带超时的上下文用于读取请求
		readCtx, readCancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
		defer readCancel()

		// 替换请求的上下文（影响请求读取）
		c.Request = c.Request.WithContext(readCtx)

		c.Next()
		//  else {
		// 	// 对其他api进行超时设置
		// 	if applyReadTimeout(c, 20*time.Second) {
		// 		return
		// 	}
		// }

		// applyWriteTimeout(c, 20*time.Second)
		// 创建带超时的上下文用于写入响应
		writeCtx, writeCancel := context.WithTimeout(c.Request.Context(), 20*time.Second)
		defer writeCancel()

		// 替换请求的上下文（影响响应写入）
		c.Request = c.Request.WithContext(writeCtx)
	}
}
