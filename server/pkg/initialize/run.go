package initialize

import (
	"server/pkg/api/router"
	"server/pkg/db"
	"server/pkg/global"
	"server/pkg/initialize/internal"
	"server/pkg/model/system"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type InitParams struct {
	Routes         []router.Router
	DatabaseTables []any
}

func RunSystem(params InitParams) {
	// 加载环境变量
	godotenv.Load()
	// 初始化Viper
	global.MPA_VP = Viper()
	// 初始化zap日志库
	global.MPA_LOG = Zap()
	zap.ReplaceGlobals(global.MPA_LOG)

	db.InitGorm() // gorm连接数据库并注册到单例

	// 注册数据库表
	registerTables(params.DatabaseTables...)

	engine := gin.New()
	// 添加全局路由前缀
	router := engine.Group(global.MPA_CONFIG.System.RouterPrefix)

	// 安装中间件
	installMiddleware(router)
	// 安装路由
	installRouter(router, params.Routes...)

	global.MPA_ROUTERS = engine.Routes()

	// 从数据库中读取所有被忽略的API
	global.MPA_DB.Model(&system.SysIgnoreApi{}).Find(&global.MPA_IGNORED_APIS)

	if global.MPA_CONFIG.System.Env == "public" {
		gin.SetMode(gin.ReleaseMode) //DebugMode ReleaseMode TestMode
	}

	if global.MPA_CONFIG.System.UseMultipoint || global.MPA_CONFIG.System.UseRedis {
		// 开启了多点登录限制，并且使用 Redis
		// 初始化redis服务
		Redis()
	}

	internal.CheckAndInitDB()

	// 启动服务器
	run(engine)
}
