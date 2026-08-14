package router

import (
	v1 "server/pkg/api/router/v1"

	"github.com/gin-gonic/gin"
)

type Router interface {
	InitRouter(Router gin.IRouter, middleware ...gin.HandlerFunc)
}

// AppRouter 应用路由结构体
// 负责初始化应用的路由配置
// 作为路由初始化的入口点
type SysRouter struct{}

// InitRouter 初始化应用路由
// 参数：
//
//	r: Gin 路由器实例
//	middleware: 全局中间件列表
//
// 功能：
//
//	初始化所有 API 路由
//	传入全局中间件，确保所有路由都能使用这些中间件
func (SysRouter) InitRouter(Router gin.IRouter, middleware ...gin.HandlerFunc) {
	v1.InitRouter(Router, middleware...)

	// 后续可以添加其他版本的路由初始化
	// v2.InitRouter(r, middleware...)
}
