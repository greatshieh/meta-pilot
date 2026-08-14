package initialize

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"server/pkg/global"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func run(engine *gin.Engine) {
	name := global.MPA_CONFIG.System.Name
	if name == "" {
		name = "MetaPilot"
	}

	msg := `
---------------------------------------------------------------------------------
  __  __   _____   _____      _              ____    ___   _        ___    _____ 
 |  \/  | | ____| |_   _|    / \            |  _ \  |_ _| | |      / _ \  |_   _|
 | |\/| | |  _|     | |     / _ \    _____  | |_) |  | |  | |     | | | |   | |  
 | |  | | | |___    | |    / ___ \  |_____| |  __/   | |  | |___  | |_| |   | |  
 |_|  |_| |_____|   |_|   /_/   \_\         |_|     |___| |_____|  \___/    |_|  



            %s server is running on address: %s
---------------------------------------------------------------------------------
`
	address := fmt.Sprintf("0.0.0.0:%d", global.MPA_CONFIG.System.Addr)

	// 关闭超时时间，等待已有请求处理完成
	WaitTime := 5 * time.Second

	// 创建http服务器实例
	srv := &http.Server{
		Addr:    address,
		Handler: engine,
	}

	// 3. 启动服务器（goroutine 异步启动，避免阻塞信号监听）
	go func() {
		fmt.Printf(msg, name, address)

		// 监听服务，除非主动关闭否则会阻塞
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.MPA_LOG.Error(name+" server listen failed", zap.Error(err))
		}
	}()

	// 4. 监听系统关闭信号（关闭核心）
	quit := make(chan os.Signal, 1)
	// 监听：Ctrl+C(SIGINT)、容器停止(SIGTERM)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSEGV)
	// 阻塞等待信号
	<-quit
	global.MPA_LOG.Warn("server is shutting down...")

	// 5. 优雅关闭服务器（设置超时上下文，防止无限等待）
	ctx, cancel := context.WithTimeout(context.Background(), WaitTime)
	defer cancel()

	// 关闭服务器：停止接收新请求，等待已有请求处理完成
	if err := srv.Shutdown(ctx); err != nil {
		global.MPA_LOG.Warn("server forced to shutdown: %v", zap.Error(err))
	}

	global.MPA_LOG.Info("server exited gracefully")
}
