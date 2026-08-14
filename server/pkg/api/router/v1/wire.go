//go:build wireinject
// +build wireinject

package v1

import (
	"server/pkg/api/handler"
	v1Handler "server/pkg/api/handler/v1"
	v1Dao "server/pkg/dao/v1"
	"server/pkg/db"
	v1Service "server/pkg/service/v1"

	"github.com/google/wire"
)

// 数据库连接 Provider
var DBProvider = wire.NewSet(
	db.GetDB,            // 获取单例数据库连接
	db.NewGormTxManager, // 创建事务管理器
)

// newApiHandler 创建 API 处理器实例
func newApiHandler() handler.ApiHandlerInterface {
	wire.Build(v1Handler.NewApiHandler, v1Service.NewApiService, v1Dao.NewApiDao, v1Dao.NewCasbinDao, v1Dao.NewAuthorityDao, DBProvider)
	return nil
}

// newAuthorityHandler 创建权限处理器实例
func newAuthorityHandler() handler.AuthorityHandlerInterface {
	wire.Build(
		v1Handler.NewAuthorityHandler,
		v1Service.NewAuthorityService,
		v1Dao.NewAuthorityDao,
		v1Dao.NewMenuDao,
		v1Dao.NewCasbinDao,
		DBProvider,
	)
	return nil
}

// newCasbinHandler 创建策略处理器实例
func newCasbinHandler() handler.CasbinHandlerInterface {
	wire.Build(v1Handler.NewCasbinHandler, v1Service.NewCasbinService, v1Dao.NewCasbinDao, v1Dao.NewApiDao, v1Dao.NewAuthorityDao, DBProvider)
	return nil
}

// newJwtHandler 创建 JWT 处理器实例
func newJwtHandler() handler.JwtHandlerInterface {
	wire.Build(v1Handler.NewJwtHandler, v1Service.NewJwtService, v1Dao.NewJwtDao, v1Dao.NewUserDao, DBProvider)
	return nil
}

// newMenuHandler 创建菜单处理器实例
func newMenuHandler() handler.MenuHandlerInterface {
	wire.Build(v1Handler.NewMenuHandler, v1Service.NewMenuService, v1Dao.NewMenuDao, v1Dao.NewAuthorityDao, DBProvider)
	return nil
}

// newUserHandler 创建用户处理器实例
func newUserHandler() handler.UserHandlerInterface {
	wire.Build(v1Handler.NewUserHandler, v1Service.NewUserService, v1Dao.NewUserDao, v1Dao.NewJwtDao, v1Service.NewJwtService, DBProvider)
	return nil
}
