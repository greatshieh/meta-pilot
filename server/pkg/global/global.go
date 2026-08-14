package global

import (
	"server/pkg/config"
	"sync"

	"go.uber.org/zap"
	"golang.org/x/sync/singleflight"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/songzhibin97/gkit/cache/local_cache"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	MPA_DB                  *gorm.DB
	MPA_DBList              map[string]*gorm.DB
	MPA_REDIS               *redis.Client
	MPA_CONFIG              config.Server
	MPA_VP                  *viper.Viper
	MPA_LOG                 *zap.Logger
	MPA_ROUTERS             gin.RoutesInfo
	MPA_IGNORED_APIS        gin.RoutesInfo
	MPA_Concurrency_Control = &singleflight.Group{}
	BlackCache              local_cache.Cache
	lock                    sync.RWMutex
)

// GetGlobalDBByDBName 通过名称获取db list中的db
func GetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	return MPA_DBList[dbname]
}

// MustGetGlobalDBByDBName 通过名称获取db 如果不存在则panic
func MustGetGlobalDBByDBName(dbname string) *gorm.DB {
	lock.RLock()
	defer lock.RUnlock()
	db, ok := MPA_DBList[dbname]
	if !ok || db == nil {
		panic("db no init")
	}
	return db
}
