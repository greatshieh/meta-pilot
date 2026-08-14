package v1

import (
	"context"
	sysReq "server/pkg/api/request"
	"server/pkg/errcode"
	"server/pkg/global"
	"server/pkg/model/common/request"
	"server/pkg/model/system"
	"server/pkg/utils"

	"github.com/marmotedu/errors"
	"gorm.io/gorm"
)

// CreateAPI 创建API
// 参数：
//
//	api: API信息
//
// 返回值：
//
//	err: 错误信息
func (a *ApiDao) CreateAPI(ctx context.Context, api system.SysApi) error {
	if err := a.GetDB(ctx).Where("path = ? AND method = ?", api.Path, api.Method).First(&system.SysApi{}).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(errcode.ErrApiAlreadyExist, "")
		}

		return errors.WithCode(errcode.ErrApiCreateFailed, "")
	}

	if err := a.GetDB(ctx).Create(&api).Error; err != nil {
		return errors.WithCode(errcode.ErrApiCreateFailed, "")
	}

	return nil
}

// DeleteAPI 删除API
// 参数：
//
//	id: API ID
//
// 返回值：
//
//	err: 错误信息
func (a *ApiDao) DeleteAPI(ctx context.Context, id uint) (system.SysApi, error) {
	var entity system.SysApi

	// 检查API是否存在
	if err := a.GetDB(ctx).First(&entity, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // api记录不存在
			return entity, errors.WithCode(errcode.ErrApiNotFound, "%s", err.Error())
		}

		return entity, errors.WithCode(errcode.ErrApiDeleteFailed, "%s", err.Error())
	}

	// 删除API
	if err := a.GetDB(ctx).Delete(&entity).Error; err != nil {
		return entity, errors.WithCode(errcode.ErrApiDeleteFailed, "%s", err.Error())
	}

	return entity, nil
}

// GetAPIInfoList 获取API信息列表
// 参数：
//
//	api: API查询条件
//	pageInfo: 分页信息
//
// 返回值：
//
//	[]system.SysApi: API列表
//	total: 总记录数
//	err: 错误信息
func (a *ApiDao) GetAPIInfoList(ctx context.Context, api system.SysApi, pageInfo request.PageInfo) ([]system.SysApi, int64, error) {
	var list []system.SysApi
	var total int64

	db := a.GetDB(ctx).Model(&system.SysApi{})

	if api.Path != "" {
		db = db.Where("path LIKE ?", "%"+api.Path+"%")
	}

	if api.ApiGroup != "" {
		db = db.Where("api_group = ?", api.ApiGroup)
	}

	if api.Method != "" {
		db = db.Where("method = ?", api.Method)
	}

	if api.Description != "" {
		db = db.Where("description LIKE ?", "%"+api.Description+"%")
	}

	if err := db.Count(&total).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.WithCode(errcode.ErrDatabase, "")
		}

		return nil, 0, nil
	}

	if err := db.Scopes(utils.Paginate(pageInfo)).Find(&list).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, errors.WithCode(errcode.ErrDatabase, "")
		}
		return nil, 0, nil
	}

	return list, total, nil
}

// UpdateAPI 更新API
// 参数：
//
//	api: API信息
//
// 返回值：
//
//	err: 错误信息
func (a *ApiDao) UpdateAPI(ctx context.Context, api system.SysApi) error {
	var oldA system.SysApi

	err := a.GetDB(ctx).First(&oldA, "id = ?", api.ID).Error
	if oldA.Path != api.Path || oldA.Method != api.Method {
		var duplicateApi system.SysApi
		if ferr := a.GetDB(ctx).First(&duplicateApi, "path = ? AND method = ?", api.Path, api.Method).Error; ferr != nil {
			if !errors.Is(ferr, gorm.ErrRecordNotFound) {
				return errors.WithCode(errcode.ErrApiUpdateFailed, "")
			}
		} else {
			if duplicateApi.ID != api.ID {
				return errors.WithCode(errcode.ErrApiAlreadyExist, "")
			}
		}
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(errcode.ErrApiAlreadyExist, "")
		}
		return errors.WithCode(errcode.ErrApiUpdateFailed, "")
	}

	// err = casbinDao.UpdateCasbinApi(oldA.Path, api.Path, oldA.Method, api.Method)
	// if err != nil {
	// 	return err
	// }

	if err := a.GetDB(ctx).Save(&api).Error; err != nil {
		return errors.WithCode(errcode.ErrApiUpdateFailed, "")
	}

	return nil
}

// GetAllAPIs 获取所有API
// 参数：
//
//	authorityID: 权限ID
//
// 返回值：
//
//	apis: API列表
//	err: 错误信息
func (a *ApiDao) GetAllAPIs(ctx context.Context, authorityID uint) ([]system.SysApi, error) {
	var apis []system.SysApi
	err := a.GetDB(ctx).Order("id desc").Find(&apis).Error
	return apis, err
}

// GetAPIGroups 获取API分组
// 返回值：
//
//	groups: API分组列表
//	err: 错误信息
func (a *ApiDao) GetAPIGroups(ctx context.Context) ([]string, error) {
	var apis []system.SysApi

	if err := a.GetDB(ctx).Find(&apis).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.WithCode(errcode.ErrApiNotFound, "")
		}
		return nil, errors.WithCode(errcode.ErrApiGroupFailed, "")
	}

	var groups = make([]string, 0)
	for i := range apis {
		newGroup := true
		for i2 := range groups {
			if groups[i2] == apis[i].ApiGroup {
				newGroup = false
			}
		}
		if newGroup {
			groups = append(groups, apis[i].ApiGroup)
		}
	}

	return groups, nil
}

// SyncAPI 同步API
// 返回值：
//
//	newApis: 新增的API列表
//	deleteApis: 删除的API列表
//	ignoreApis: 忽略的API列表
//	err: 错误信息
func (a *ApiDao) SyncAPI(ctx context.Context) (newApis, deleteApis, ignoreApis []system.SysApi, err error) {
	newApis = make([]system.SysApi, 0)
	deleteApis = make([]system.SysApi, 0)
	ignoreApis = make([]system.SysApi, 0)

	var apis []system.SysApi

	if err = a.GetDB(ctx).Find(&apis).Error; err != nil {
		return nil, nil, nil, errors.WithCode(errcode.ErrApiRegisterListFailed, "")
	}

	var ignores []system.SysIgnoreApi
	if err = a.GetDB(ctx).Find(&ignores).Error; err != nil {
		return nil, nil, nil, errors.WithCode(errcode.ErrApiRegisterListFailed, "")
	}

	for i := range ignores {
		ignoreApis = append(ignoreApis, system.SysApi{
			Path:        ignores[i].Path,
			Description: "",
			ApiGroup:    "",
			Method:      ignores[i].Method,
		})
	}

	var cacheApis []system.SysApi
	for i := range global.MPA_ROUTERS {
		ignoresFlag := false
		for j := range ignores {
			// 标记被忽略的api
			if ignores[j].Path == global.MPA_ROUTERS[i].Path && ignores[j].Method == global.MPA_ROUTERS[i].Method {
				ignoresFlag = true
			}
		}

		// api没有被忽略，加入缓存数组
		if !ignoresFlag {
			cacheApis = append(cacheApis, system.SysApi{
				Path:   global.MPA_ROUTERS[i].Path,
				Method: global.MPA_ROUTERS[i].Method,
			})
		}
	}

	//如果内存中的api不存在于数据库中，则把api放入新增数组
	for i := range cacheApis {
		var flag bool

		for j := range apis {
			if cacheApis[i].Path == apis[j].Path && cacheApis[i].Method == apis[j].Method {
				// 标记同时存在于内存和数据库中的api
				flag = true
				break
			}
		}

		// 如果存在于内存但不在api数组中，说明是新增的api
		if !flag {
			newApis = append(newApis, system.SysApi{
				Path:        cacheApis[i].Path,
				Description: "",
				ApiGroup:    "",
				Method:      cacheApis[i].Method,
			})
		}
	}

	// 如果数据库中的api不在内存中，则把api放入删除数组，
	for i := range apis {
		var flag bool
		// 标记同时存在于数据库和内存中的api
		for j := range cacheApis {
			// 同时存在api数组和内存，说明是正常的api，标记
			if cacheApis[j].Path == apis[i].Path && cacheApis[j].Method == apis[i].Method {
				flag = true
				break
			}
		}
		// 没有被标记，说明是删除的api
		if !flag {
			deleteApis = append(deleteApis, apis[i])
		}
	}
	return
}

// IgnoreAPI 忽略API
// 参数：
//
//	ignoreApi: 忽略的API信息
//
// 返回值：
//
//	err: 错误信息
func (a *ApiDao) IgnoreAPI(ctx context.Context, ignoreApi system.SysIgnoreApi) (err error) {
	if ignoreApi.Flag {
		if err := a.GetDB(ctx).Create(&ignoreApi).Error; err != nil {
			return errors.WithCode(errcode.ErrApiIgnoreFailed, "")
		}
		return nil
	}
	return a.GetDB(ctx).Unscoped().Delete(&ignoreApi, "path = ? AND method = ?", ignoreApi.Path, ignoreApi.Method).Error
}

// EnterSyncAPI 进入同步API
// 参数：
//
//	syncApis: 同步API参数
//
// 返回值：
//
//	err: 错误信息
func (a *ApiDao) EnterSyncAPI(ctx context.Context, syncApis sysReq.EnterSyncApiParams) (err error) {
	var txErr error
	if len(syncApis.NewApis) > 0 {
		txErr = a.GetDB(ctx).Create(&syncApis.NewApis).Error
		if txErr != nil {
			return txErr
		}
	}

	for i := range syncApis.DeleteApis {
		if txErr := a.GetDB(ctx).Delete(&system.SysApi{}, "path = ? AND method = ?", syncApis.DeleteApis[i].Path, syncApis.DeleteApis[i].Method).Error; txErr != nil {
			return txErr
		}
	}

	return nil
}
