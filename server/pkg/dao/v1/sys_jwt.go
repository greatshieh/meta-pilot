package v1

import (
	"context"
	"server/pkg/errcode"
	"server/pkg/utils"

	"server/pkg/global"
	"server/pkg/model/system"

	"github.com/marmotedu/errors"
	"go.uber.org/zap"
)

func (jwtDao *JwtDao) JsonInBlacklist(ctx context.Context, jwtList system.JwtBlacklist) error {
	if err := jwtDao.GetDB(ctx).Create(&jwtList).Error; err != nil {
		return errors.WithCode(errcode.ErrTokenRevokeFailed, "%s", err.Error())
	}

	global.BlackCache.SetDefault(jwtList.Jwt, struct{}{})

	return nil
}

func (jwtDao *JwtDao) IsBlacklist(ctx context.Context, jwt string) bool {
	_, ok := global.BlackCache.Get(jwt)
	return ok
	// err := global.MPA_DB.Where("jwt = ?", jwt).First(&system.JwtBlacklist{}).Error
	// isNotFound := errors.Is(err, gorm.ErrRecordNotFound)
	// return !isNotFound
}

func (jwtDao *JwtDao) GetRedisJWT(ctx context.Context, userName string) (redisJWT string, err error) {
	redisJWT, err = global.MPA_REDIS.Get(ctx, userName).Result()
	return redisJWT, err
}

func (jwtDao *JwtDao) SetRedisJWT(ctx context.Context, jwt, userName string) (err error) {
	// 此处过期时间等于jwt过期时间
	dr, err := utils.ParseDuration(global.MPA_CONFIG.JWT.ExpiresTime)
	if err != nil {
		return errors.WithCode(errcode.ErrTokenRedisStateSetFailed, "%s", err.Error())
	}
	timer := dr
	if err = global.MPA_REDIS.Set(ctx, userName, jwt, timer).Err(); err != nil {
		return errors.WithCode(errcode.ErrTokenRedisStateSetFailed, "%s", err.Error())
	}
	return nil
}

func (jwtDao *JwtDao) LoadAll(ctx context.Context) {
	var data []string
	err := jwtDao.GetDB(ctx).Model(&system.JwtBlacklist{}).Select("jwt").Find(&data).Error
	if err != nil {
		global.MPA_LOG.Error("加载数据库jwt黑名单失败!", zap.Error(err))
		return
	}
	for i := 0; i < len(data); i++ {
		global.BlackCache.SetDefault(data[i], struct{}{})
	} // jwt黑名单 加入 BlackCache 中
}
