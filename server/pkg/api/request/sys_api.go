package request

import (
	"server/pkg/model/common/request"
	"server/pkg/model/system"
)

type SearchApiParams struct {
	system.SysApi
	request.PageInfo
}

type EnterSyncApiParams struct {
	NewApis    []system.SysApi `json:"newApis"`    // 新增api数组
	DeleteApis []system.SysApi `json:"deleteApis"` // 删除api数组
}
