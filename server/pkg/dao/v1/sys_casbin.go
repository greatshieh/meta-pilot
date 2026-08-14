package v1

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/utils"

	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/marmotedu/errors"
)

// SetCasbin 设置Casbin权限
// 参数：
//
//	adminAuthorityID: 管理员权限ID
//	authorityID: 目标权限ID
//	casbinInfos: Casbin权限信息列表
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	为指定权限设置Casbin权限规则，包括权限验证、去重处理和规则添加
func (casbinDao *CasbinDao) SetCasbin(ctx context.Context, rules [][]string) error {
	e := utils.GetCasbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.WithCode(errcode.ErrCasbinSet, "存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// UpdateCasbinApi API更新随动
// 参数：
//
//	oldPath: 旧API路径
//	newPath: 新API路径
//	oldMethod: 旧HTTP方法
//	newMethod: 新HTTP方法
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	当API路径或方法变更时，同步更新Casbin权限规则
func (casbinDao *CasbinDao) UpdateCasbinApi(ctx context.Context, oldPath, newPath, oldMethod, newMethod string) error {
	err := global.MPA_DB.Model(&gormadapter.CasbinRule{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]any{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	if err != nil {
		return err
	}

	e := utils.GetCasbin()
	return e.LoadPolicy()
}

// GetPolicyPathByAuthorityID 根据权限ID获取权限列表
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	pathMaps: Casbin权限信息列表
//
// 功能：
//
//	获取指定权限ID的所有Casbin权限规则
func (casbinDao *CasbinDao) GetPolicyPathByAuthorityID(ctx context.Context, authorityID string) (pathMaps []request.CasbinInfo) {
	e := utils.GetCasbin()
	list, _ := e.GetFilteredPolicy(0, authorityID)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

// ClearCasbin 清除匹配的权限
// 参数：
//
//	v: 过滤索引
//	p: 过滤参数
//
// 返回值：
//
//	bool: 操作是否成功
//
// 功能：
//
//	清除指定过滤条件的Casbin权限规则
func (casbinDao *CasbinDao) ClearCasbin(ctx context.Context, v int, p ...string) (bool, error) {
	e := utils.GetCasbin()
	return e.RemoveFilteredPolicy(v, p...)
}

// RemoveFilteredPolicy 使用数据库方法清理筛选的权限
// 参数：
//
//	db: GORM数据库连接
//	authorityID: 权限ID
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	使用数据库方法删除指定权限ID的Casbin规则
//	此方法需要调用FreshCasbin方法才可以在系统中即刻生效
func (casbinDao *CasbinDao) RemoveFilteredPolicy(ctx context.Context, authorityID string) error {
	return casbinDao.GetDB(ctx).Delete(&gormadapter.CasbinRule{}, "v0 = ?", authorityID).Error
}

// SyncPolicy 同步目前数据库的权限
// 参数：
//
//	db: GORM数据库连接
//	authorityID: 权限ID
//	rules: 权限规则列表
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	同步指定权限ID的Casbin规则到数据库
//	此方法需要调用FreshCasbin方法才可以在系统中即刻生效
func (casbinDao *CasbinDao) SyncPolicy(ctx context.Context, authorityID string, rules [][]string) error {
	err := casbinDao.RemoveFilteredPolicy(ctx, authorityID)
	if err != nil {
		return err
	}
	return casbinDao.AddPolicies(ctx, rules)
}

// AddPolicies 添加匹配的权限
// 参数：
//
//	db: GORM数据库连接
//	rules: 权限规则列表
//
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	批量添加Casbin权限规则到数据库
func (casbinDao *CasbinDao) AddPolicies(ctx context.Context, rules [][]string) error {
	var casbinRules []gormadapter.CasbinRule
	for i := range rules {
		casbinRules = append(casbinRules, gormadapter.CasbinRule{
			Ptype: "p",
			V0:    rules[i][0],
			V1:    rules[i][1],
			V2:    rules[i][2],
		})
	}
	return casbinDao.GetDB(ctx).Create(&casbinRules).Error
}

// FreshCasbin 刷新Casbin策略
// 返回值：
//
//	err: 错误信息
//
// 功能：
//
//	重新加载Casbin策略，使数据库中的更改立即生效
func (casbinDao *CasbinDao) FreshCasbin(ctx context.Context) (err error) {
	e := utils.GetCasbin()
	err = e.LoadPolicy()
	return err
}
