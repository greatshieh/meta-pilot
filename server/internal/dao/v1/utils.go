package v1

import (
	"gorm.io/gorm"
)

// 定义一个全局的通道，用于外部触发消息推送
var SSEEventChan = make(chan string)

// whereGrade 根据年级过滤 GORM 查询
// 参数:
//   - orm: GORM 查询对象
//   - grade: 年级字符串
//
// 返回值:
//   - *gorm.DB: 过滤后的 GORM 查询对象
func whereGrade(orm *gorm.DB, table, grade string) *gorm.DB {
	// 检查年级是否为空
	if grade != "" {
		// 如果年级不为空，添加 WHERE 条件到 GORM 查询对象
		orm = orm.Where(table+".grade = ?", grade)
	}

	// 返回过滤后的 GORM 查询对象
	return orm
}

// whereClass 根据班级过滤 GORM 查询
// 参数:
//   - orm: GORM 查询对象
//   - class: 班级字符串
//
// 返回值:
//   - *gorm.DB: 过滤后的 GORM 查询对象
func whereClass(orm *gorm.DB, table, class string) *gorm.DB {
	// 检查班级是否为空
	if class != "" {
		// 如果班级不为空，添加 WHERE 条件到 GORM 查询对象
		orm = orm.Where(table+".class = ?", class)
	}

	// 返回过滤后的 GORM 查询对象
	return orm
}

// whereName 根据姓名过滤 GORM 查询
// 参数:
//   - orm: GORM 查询对象
//   - name: 姓名字符串
//
// 返回值:
//   - *gorm.DB: 过滤后的 GORM 查询对象
func whereName(orm *gorm.DB, table, name string) *gorm.DB {
	// 检查姓名是否为空
	if name != "" {
		// 如果姓名不为空，添加 WHERE 条件到 GORM 查询对象，使用 LIKE 进行模糊查询
		switch table {
		case "student_info":
			orm = orm.Where(table+".name LIKE ?", "%"+name+"%")
		case "examination_models":
			orm = orm.Where(table+".examName LIKE ?", "%"+name+"%")
		}

	}

	// 返回过滤后的 GORM 查询对象
	return orm
}

// scopesBase 生成一个 GORM 查询作用域，用于根据年级、班级和姓名过滤学生信息
// 参数:
//   - grade: 年级字符串
//   - class: 班级字符串
//   - name: 姓名字符串
//
// 返回值:
//   - func(db *gorm.DB) *gorm.DB: 一个 GORM 查询作用域函数，用于链式调用
func scopesBase(table string, grade, class, name string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// 链式调用 whereGrade、whereClass 和 whereName 函数，根据条件过滤 GORM 查询对象
		return whereGrade(whereClass(whereName(db, table, name), table, class), table, grade)
	}
}
