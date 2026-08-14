package system

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/global"
	"strconv"

	"github.com/marmotedu/errors"
)

// SetCasbin 设置 Casbin 策略
// 参数：
//
//	adminAuthorityID: 管理员权限 ID
//	authorityID: 权限 ID
//	cmr: Casbin 策略信息
//
// 返回值：
//
//	error: 错误信息
func (cas *CasbinService) SetCasbin(ctx context.Context, adminAuthorityID, authorityID uint, casbinInfos []request.CasbinInfo) error {
	// 检查管理员权限是否足够
	if err := cas.authorityDao.CheckAuthorityIDAuth(ctx, adminAuthorityID, authorityID); err != nil {
		return errors.WithCode(errcode.ErrCasbinSet, err.Error())
	}

	if global.MPA_CONFIG.System.UseStrictAuth {
		apis, err := cas.apiDao.GetAllAPIs(ctx, adminAuthorityID)
		if err != nil {
			return errors.WithCode(errcode.ErrCasbinSet, err.Error())
		}

		for i := range casbinInfos {
			hasApi := false
			for j := range apis {
				if apis[j].Path == casbinInfos[i].Path && apis[j].Method == casbinInfos[i].Method {
					hasApi = true
					break
				}
			}
			if !hasApi {
				return errors.WithCode(errcode.ErrCasbinSet, "存在api不在权限列表中")
			}
		}
	}

	authorityIDStr := strconv.Itoa(int(authorityID))
	cas.casbinDao.ClearCasbin(ctx, 0, authorityIDStr)
	rules := [][]string{}
	// 做权限去重处理
	deduplicateMap := make(map[string]bool)
	for _, v := range casbinInfos {
		key := authorityIDStr + v.Path + v.Method
		if _, ok := deduplicateMap[key]; !ok {
			deduplicateMap[key] = true
			rules = append(rules, []string{authorityIDStr, v.Path, v.Method})
		}
	}
	if len(rules) == 0 {
		return nil
	} // 设置空权限无需调用 AddPolicies 方法

	err := cas.casbinDao.SetCasbin(ctx, rules)

	return err
}

// GetPolicyPathByAuthorityId 根据权限 ID 获取策略路径
// 参数：
//
//	authorityID: 权限 ID
//
// 返回值：
//
//	[]request.CasbinInfo: Casbin 策略信息列表
func (cas *CasbinService) GetPolicyPathByAuthorityID(ctx context.Context, authorityID string) []request.CasbinInfo {
	return cas.casbinDao.GetPolicyPathByAuthorityID(ctx, authorityID)
}
