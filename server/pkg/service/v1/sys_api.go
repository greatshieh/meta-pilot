package system

import (
	"context"
	"server/pkg/api/request"
	"server/pkg/global"
	"server/pkg/model/common/response"
	"server/pkg/model/system"
	"strconv"
)

// CreateAPI 创建 API
// 参数：
//
//	api: 要创建的 API 信息
//
// 返回值：
//
//	system.SysApi: 创建的 API 信息
//	error: 错误信息
func (a *ApiService) CreateAPI(ctx context.Context, api system.SysApi) (system.SysApi, error) {
	return api, a.apiDao.CreateAPI(ctx, api)
}

// GetAPIList 获取 API 列表
// 参数：
//
//	pageInfo: 分页和搜索参数
//
// 返回值：
//
//	response.PageResult: 分页结果，包含 API 列表和分页信息
//	error: 错误信息
func (a *ApiService) GetAPIList(ctx context.Context, pageInfo request.SearchApiParams) (response.PageResult, error) {
	list, total, err := a.apiDao.GetAPIInfoList(ctx, pageInfo.SysApi, pageInfo.PageInfo)
	if err != nil {
		return response.PageResult{}, err
	}

	return response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, nil
}

// GetAllAPIs 获取所有 API
// 参数：
//
//	authorityID: 角色 ID
//
// 返回值：
//
//	[]system.SysApi: API 列表
//	error: 错误信息
func (a *ApiService) GetAllAPIs(ctx context.Context, authorityID uint) ([]system.SysApi, error) {
	var finalAPIs []system.SysApi
	err := a.tx.Transaction(func(ctx context.Context) error {
		// 查找当前权限ID的父权限
		parentAuthorityID, err := a.authorityDao.GetParentAuthorityID(ctx, authorityID)
		if err != nil {
			return err
		}

		apis, err := a.apiDao.GetAllAPIs(ctx, authorityID)
		// 如果是一级权限，或者不是严格权限策略，直接返回
		if parentAuthorityID == 0 || !global.MPA_CONFIG.System.UseStrictAuth {
			finalAPIs = apis
			return err
		}

		authorityIDStr := strconv.Itoa(int(authorityID))
		paths := a.casbinDao.GetPolicyPathByAuthorityID(ctx, authorityIDStr)
		// 挑选 apis里面的path和method也在paths里面的api
		var authAPIs []system.SysApi
		for i := range apis {
			for j := range paths {
				if paths[j].Path == apis[i].Path && paths[j].Method == apis[i].Method {
					authAPIs = append(authAPIs, apis[i])
				}
			}
		}

		finalAPIs = authAPIs
		return err
	})

	return finalAPIs, err
}

// DeleteAPI 删除 API
// 参数：
//
//	reqID: 要删除的 API ID
//
// 返回值：
//
//	error: 错误信息
func (a *ApiService) DeleteAPI(ctx context.Context, reqID uint) error {
	// todo 管理事物
	if entity, err := a.apiDao.DeleteAPI(ctx, reqID); err != nil {
		return err
	} else {
		a.casbinDao.ClearCasbin(ctx, 1, entity.Path, entity.Method)
	}

	return nil
}

// UpdateAPI 更新 API
// 参数：
//
//	api: 要更新的 API 信息
//
// 返回值：
//
//	system.SysApi: 更新后的 API 信息
//	error: 错误信息
func (a *ApiService) UpdateAPI(ctx context.Context, api system.SysApi) (system.SysApi, error) {
	return api, a.apiDao.UpdateAPI(ctx, api)
}

// GetAPIGroups 获取 API 分组
// 返回值：
//
//	[]string: API 分组列表
//	error: 错误信息
func (a *ApiService) GetAPIGroups(ctx context.Context) (groups []string, err error) {
	return a.apiDao.GetAPIGroups(ctx)
}

// SyncAPI 同步 API
// 返回值：
//
//	[]system.SysApi: 新增的 API 列表
//	[]system.SysApi: 更新的 API 列表
//	[]system.SysApi: 删除的 API 列表
//	error: 错误信息
func (a *ApiService) SyncAPI(ctx context.Context) (newApis, deleteApis, ignoreApis []system.SysApi, err error) {
	return a.apiDao.SyncAPI(ctx)
}

// IgnoreAPI 忽略 API
// 参数：
//
//	ignoreApi: 要忽略的 API 信息
//
// 返回值：
//
//	error: 错误信息
func (a *ApiService) IgnoreAPI(ctx context.Context, ignoreApi system.SysIgnoreApi) error {
	return a.apiDao.IgnoreAPI(ctx, ignoreApi)
}

// EnterSyncAPI 进入同步 API
// 参数：
//
//	enterSync: 同步参数
//
// 返回值：
//
//	error: 错误信息
func (a *ApiService) EnterSyncAPI(ctx context.Context, enterSync request.EnterSyncApiParams) error {
	return a.tx.Transaction(func(ctx context.Context) error {
		// 从casbin中删除的API
		for i := range enterSync.DeleteApis {
			if _, err := a.casbinDao.ClearCasbin(ctx, 1, enterSync.DeleteApis[i].Path, enterSync.DeleteApis[i].Method); err != nil {
				return err
			}
		}

		// 在数据库同步api
		if err := a.apiDao.EnterSyncAPI(ctx, enterSync); err != nil {
			return err
		}

		return nil
	})
}
