// 服务入口
// 控制服务的版本注册
package service

import v1 "server/internal/service/v1"

type ServiceGroup struct {
}

var WecharService = new(v1.WechatService)

// var SomeService = new(v2.SomeService)
