package utils

import (
	"server/pkg/model/common/request"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Paginate(pageInfo request.PageInfo) func(db *gorm.DB) *gorm.DB {
	// 计算分页偏移量
	offset := pageInfo.PageSize * (pageInfo.Page - 1)
	return func(db *gorm.DB) *gorm.DB {
		if pageInfo.PageSize > 100 {
			pageInfo.PageSize = 100
		}

		db = db.Offset(offset).Limit(pageInfo.PageSize)

		if pageInfo.SortField != "" {
			if pageInfo.SortOrder == "desc" {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: pageInfo.SortField}, Desc: true})
			} else {
				db = db.Order(clause.OrderByColumn{Column: clause.Column{Name: pageInfo.SortField}, Desc: false})
			}
		}

		return db
	}
}
