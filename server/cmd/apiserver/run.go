package main

import (
	"context"
	"errors"
	"fmt"
	"server/internal/api/router"
	sysRouter "server/pkg/api/router"
	"server/pkg/global"
	"server/pkg/initialize"
	"server/pkg/model/system"
	"server/pkg/utils"

	uuid "github.com/satori/go.uuid"
)

func main() {
	// 生成需要初始化的参数
	// initialize.InitParams{
	// 需要注册的路由
	// Routes: router.AppRouter(),
	// 需要注册的自定义数据表
	// DatabaseTables: []any{&system.SysUser{}},
	// }

	routes := new(router.AppRouter)

	params := initialize.InitParams{
		Routes: []sysRouter.Router{routes},
	}

	// 运行系统
	initialize.RunSystem(params)

	// 程序结束前关闭数据库链接
	db, _ := global.MPA_DB.DB()
	if db != nil {
		defer db.Close()
	}
}

type ctxKey string

// 初始化数据库
func InitDB() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer func(c func()) { c() }(cancel)

	var err error

	// 初始化系统角色
	if ctx, err = initSysAuthor(ctx); err != nil {
		return
	}

	// 初始化菜单
	if ctx, err = initMenu(ctx); err != nil {
		return
	}

	if ctx, err = initAuthMenu(ctx); err != nil {
		return
	}

	initUser(ctx)
}

// 初始化角色表
func initSysAuthor(ctx context.Context) (context.Context, error) {
	var superAdmin uint = 0
	entities := []system.SysRoleAuthority{
		{AuthorityID: 200, AuthorityName: "系统管理员", ParentID: &superAdmin, DefaultRouter: "dashboard"},
	}

	if err := global.MPA_DB.Create(&entities).Error; err != nil {
		fmt.Printf("%s表数据初始化失败!", system.SysRoleAuthority{}.TableName())
		return ctx, fmt.Errorf("%s数据初始化失败", system.SysRoleAuthority{}.TableName())
	}

	// 为超级管理员添加所有子角色
	if err := global.MPA_DB.Model(&entities[0]).Association("DataAuthorityId").Replace(
		[]*system.SysRoleAuthority{
			{AuthorityID: 200},
		}); err != nil {
		return ctx, fmt.Errorf("%s数据初始化失败", system.SysRoleAuthority{}.TableName())
	}

	ctx = context.WithValue(ctx, ctxKey("auth"), entities)

	return ctx, nil
}

// 初始化菜单
func initMenu(ctx context.Context) (context.Context, error) {

	// 定义所有菜单
	allMenus := []system.SysMenu{
		{MenuLevel: 0, ParentID: 0, Path: "dashboard", Name: "dashboard", Component: "views/dashboard/index.vue", Sort: 1, Meta: system.Meta{Title: "首页", Icon: "ant-design:dashboard-outlined", Affix: true}},
		{MenuLevel: 0, ParentID: 0, Path: "student", Name: "student", Component: "views/student/index.vue", Sort: 2, Meta: system.Meta{Title: "学生管理", Icon: "formkit:people"}},
		{MenuLevel: 0, ParentID: 0, Path: "score", Name: "score", Component: "views/score/index.vue", Sort: 3, Meta: system.Meta{Title: "成绩管理", Icon: "ant-design:line-chart-outlined"}},
		{MenuLevel: 0, ParentID: 0, Path: "exam", Name: "exam", Component: "views/exambank/index.vue", Sort: 4, Meta: system.Meta{Title: "考试管理", Icon: "ant-design:line-chart-outlined"}},
		{MenuLevel: 0, ParentID: 0, Path: "works", Name: "works", Component: "views/works/index.vue", Sort: 5, Meta: system.Meta{Title: "作品管理", Icon: "fluent-mdl2:pen-workspace"}},
		{MenuLevel: 0, ParentID: 0, Path: "changpassword", Name: "changepwd", Component: "views/changpassword/index.vue", Sort: 6, Meta: system.Meta{Title: "修改密码", Hidden: true}},
		{MenuLevel: 0, ParentID: 0, Path: "sysManager", Name: "sysManager", Component: "views/sysManager/index.vue", Sort: 7, Meta: system.Meta{Title: "系统管理", Icon: "tdesign:system-setting"}},
	}

	// 先创建父级菜单（ParentId = 0 的菜单）
	if err := global.MPA_DB.Create(&allMenus).Error; err != nil {
		return ctx, fmt.Errorf("❌ 父级菜单初始化失败！")
	}

	// 建立菜单映射 - 通过Name查找已创建的菜单及其ID
	menuNameMap := make(map[string]uint)
	for _, menu := range allMenus {
		menuNameMap[menu.Name] = menu.ID
	}

	subMenus := []system.SysMenu{
		// 学生管理子菜单
		{MenuLevel: 0, ParentID: menuNameMap["student"], Path: "error", Name: "studenterr", Component: "views/student/error.vue", Sort: 1, Meta: system.Meta{Title: "错误详情", Hidden: true, ActiveName: "student"}},
		// 成绩管理子菜单
		{MenuLevel: 0, ParentID: menuNameMap["score"], Path: "grade", Name: "gradescore", Component: "views/score/grade.vue", Sort: 1, Meta: system.Meta{Title: "年级成绩"}},
		{MenuLevel: 0, ParentID: menuNameMap["score"], Path: "class", Name: "classscore", Component: "views/score/class.vue", Sort: 2, Meta: system.Meta{Title: "班级成绩"}},
		{MenuLevel: 0, ParentID: menuNameMap["score"], Path: "list", Name: "scorelist", Component: "views/score/list.vue", Sort: 3, Meta: system.Meta{Title: "成绩列表"}},
		{MenuLevel: 0, ParentID: menuNameMap["score"], Path: "upload", Name: "uploadscore", Component: "views/score/upload.vue", Sort: 4, Meta: system.Meta{Title: "上传成绩"}},
		// 试卷管理子菜单
		{MenuLevel: 0, ParentID: menuNameMap["exam"], Path: "upload", Name: "uploadexam", Component: "views/exambank/upload.vue", Sort: 1, Meta: system.Meta{Title: "上传试卷", Hidden: true, ActiveName: "exam"}},
		// 作品管理子菜单
		{MenuLevel: 0, ParentID: menuNameMap["works"], Path: "art", Name: "reviewart", Component: "views/works/review.vue", Sort: 1, Meta: system.Meta{Title: "美术作品"}},
		{MenuLevel: 0, ParentID: menuNameMap["works"], Path: "music", Name: "reviewmusic", Component: "views/works/review.vue", Sort: 2, Meta: system.Meta{Title: "音乐作品"}},
		{MenuLevel: 0, ParentID: menuNameMap["works"], Path: "upload/:subject", Name: "uploadworks", Component: "views/works/upload.vue", Sort: 3, Meta: system.Meta{Title: "上传作品", Hidden: true}},
		// 系统管理子菜单
		{MenuLevel: 0, ParentID: menuNameMap["sysManager"], Path: "authority", Name: "authManager", Component: "views/sysManager/authority/index.vue", Sort: 1, Meta: system.Meta{Title: "权限管理"}},
		{MenuLevel: 0, ParentID: menuNameMap["sysManager"], Path: "menu", Name: "menuManager", Component: "views/sysManager/menu/index.vue", Sort: 2, Meta: system.Meta{Title: "菜单管理"}},
		{MenuLevel: 0, ParentID: menuNameMap["sysManager"], Path: "sysuser", Name: "sysUser", Component: "views/sysManager/authority/index.vue", Sort: 3, Meta: system.Meta{Title: "用户管理"}},
	}

	// 创建子菜单
	if err := global.MPA_DB.Create(&subMenus).Error; err != nil {
		return ctx, fmt.Errorf("❌ 子级菜单初始化失败！")
	}

	allEntries := append(allMenus, subMenus...)
	ctx = context.WithValue(ctx, ctxKey("basemenu"), allEntries)

	return ctx, nil
}

// 关联菜单
func initAuthMenu(ctx context.Context) (context.Context, error) {
	authorities, ok := ctx.Value(ctxKey("auth")).([]system.SysRoleAuthority)
	if !ok {
		return ctx, errors.New("获取角色错误")
	}

	allMenus, ok := ctx.Value(ctxKey("basemenu")).([]system.SysMenu)
	if !ok {
		return ctx, errors.New("获取菜单错误")
	}

	// 构建菜单ID映射，方便快速查找
	menuMap := make(map[uint]system.SysMenu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
	}

	// 为不同角色分配不同权限
	// 1. 超级管理员角色(200) - 拥有所有菜单权限
	if err := global.MPA_DB.Model(&authorities[0]).Association("SysBaseMenus").Replace(allMenus); err != nil {
		return ctx, errors.New("为超级管理员分配菜单失败")
	}

	return ctx, nil
}

func initUser(ctx context.Context) (context.Context, error) {
	entities := []system.SysUser{
		{
			UUID:        uuid.NewV4(),
			UserName:    "admin",
			Password:    utils.BcryptHash("11223344"),
			NickName:    "超级管理员",
			Avatar:      "https://img.tuxiangyan.com/zb_users/upload/2023/02/202302091675904134770942.jpg",
			AuthorityID: 200,
			Phone:       "17700000000",
			Email:       "88888888@qq.com",
			IsActive:    true,
			IsStaff:     true,
		},
	}

	if err := global.MPA_DB.Create(&entities).Error; err != nil {
		return ctx, errors.New("用户表数据初始化失败")
	}

	authorityEntities, ok := ctx.Value(ctxKey("auth")).([]system.SysRoleAuthority)
	if !ok {
		return ctx, errors.New("创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err := global.MPA_DB.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return ctx, err
	}
	return ctx, nil
}
