package internal

import (
	"context"
	"errors"
	"fmt"
	"os"
	"server/pkg/global"
	"server/pkg/model/system"
	"server/pkg/utils"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type ctxKey string

// 初始化数据库
func CheckAndInitDB() {
	var initFlag system.InitFlag
	global.MPA_DB.First(&initFlag)

	if !initFlag.IsInitialized {
		// 数据库未初始化，执行初始化
		global.MPA_LOG.Info("🆘 首次部署，开始执行数据库初始化...")

		if err := initDB(); err != nil {
			global.MPA_LOG.Error("❌ 数据库初始化失败: " + err.Error())
			return
		}

		if initFlag.ID == 0 {
			initFlag.IsInitialized = true
			global.MPA_DB.Create(&system.InitFlag{IsInitialized: true})
		} else {
			global.MPA_DB.Model(&initFlag).Update("is_initialized", true)
		}

		global.MPA_LOG.Info("✅ 数据库初始化完成")
	} else {
		global.MPA_LOG.Info("✅ 数据库已初始化, 跳过")
	}

}

func initDB() error {
	return global.MPA_DB.Transaction(func(tx *gorm.DB) error {
		ctx, cancel := context.WithCancel(context.TODO())
		defer func(c func()) { c() }(cancel)

		var err error
		// 初始化系统角色
		if ctx, err = initSysRoleAuthority(ctx, tx); err != nil {
			return err
		}

		// 初始化菜单
		if ctx, err = initMenu(ctx, tx); err != nil {
			return err
		}

		if ctx, err = initAuthMenu(ctx, tx); err != nil {
			return err
		}

		if ctx, err = initUser(ctx, tx); err != nil {
			return err
		}

		return nil
	})
}

// 初始化角色表
func initSysRoleAuthority(ctx context.Context, tx *gorm.DB) (context.Context, error) {
	var superAdmin uint = 0
	entities := []system.SysRoleAuthority{
		{AuthorityID: 200, AuthorityName: "系统管理员", ParentID: &superAdmin, DefaultRouter: "dashboard"},
	}

	if err := tx.Create(&entities).Error; err != nil {
		return ctx, fmt.Errorf("%s数据初始化失败", system.SysRoleAuthority{}.TableName())
	}

	// 为超级管理员添加所有子角色
	if err := tx.Model(&entities[0]).Association("AuthorizedSubRoles").Replace(
		[]*system.SysRoleAuthority{
			{AuthorityID: 200},
		}); err != nil {
		return ctx, fmt.Errorf("%s数据初始化失败", system.SysRoleAuthority{}.TableName())
	}

	ctx = context.WithValue(ctx, ctxKey("auth"), entities)

	return ctx, nil
}

// 初始化菜单
func initMenu(ctx context.Context, tx *gorm.DB) (context.Context, error) {
	// 定义所有菜单
	allMenus := []system.SysMenu{
		{ParentID: 0, Path: "dashboard", Name: "dashboard", Component: "views/dashboard/index.vue", Sort: 1, Meta: system.Meta{Title: "首页", Icon: "ant-design:dashboard-outlined", Affix: true}},
		{ParentID: 0, Path: "changpassword", Name: "changepwd", Component: "views/changpassword/index.vue", Sort: 6, Meta: system.Meta{Title: "修改密码", Hidden: true}},
		{ParentID: 0, Path: "sysManager", Name: "sysManager", Component: "views/sysManager/index.vue", Sort: 7, Meta: system.Meta{Title: "系统管理", Icon: "tdesign:system-setting"}},
	}

	// 先创建父级菜单（ParentID = 0 的菜单）
	if err := tx.Create(&allMenus).Error; err != nil {
		return ctx, errors.New("❌ 父级菜单初始化失败！")
	}

	// 建立菜单映射 - 通过Name查找已创建的菜单及其ID
	menuNameMap := make(map[string]uint)
	for _, menu := range allMenus {
		menuNameMap[menu.Name] = menu.ID
	}

	subMenus := []system.SysMenu{
		// 系统管理子菜单
		{ParentID: menuNameMap["sysManager"], Path: "authority", Name: "authManager", Component: "views/sysManager/authority/index.vue", Sort: 1, Meta: system.Meta{Title: "权限管理"}},
		{ParentID: menuNameMap["sysManager"], Path: "menu", Name: "menuManager", Component: "views/sysManager/menu/index.vue", Sort: 2, Meta: system.Meta{Title: "菜单管理"}},
		{ParentID: menuNameMap["sysManager"], Path: "sysuser", Name: "sysUser", Component: "views/sysManager/authority/index.vue", Sort: 3, Meta: system.Meta{Title: "用户管理"}},
	}

	// 创建子菜单
	if err := tx.Create(&subMenus).Error; err != nil {
		return ctx, errors.New("❌ 子级菜单初始化失败！")
	}

	allEntries := append(allMenus, subMenus...)
	ctx = context.WithValue(ctx, ctxKey("menu"), allEntries)

	return ctx, nil
}

// 关联菜单
func initAuthMenu(ctx context.Context, tx *gorm.DB) (context.Context, error) {
	authorities, ok := ctx.Value(ctxKey("auth")).([]system.SysRoleAuthority)
	if !ok {
		return ctx, errors.New("获取角色错误")
	}

	allMenus, ok := ctx.Value(ctxKey("menu")).([]system.SysMenu)
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
	if err := tx.Model(&authorities[0]).Association("SysMenu").Replace(allMenus); err != nil {
		return ctx, errors.New("超级管理员分配菜单失败: " + err.Error())
	}

	return ctx, nil
}

func initUser(ctx context.Context, tx *gorm.DB) (context.Context, error) {
	// 从环境变量获取初始化密码
	initPassword := os.Getenv("INIT_ADMIN_PASSWORD")
	if initPassword == "" {
		return ctx, errors.New("❌ 未配置 INIT_ADMIN_PASSWORD 环境变量，禁止初始化")
	}

	var superAdmin uint = 0
	entities := []system.SysUser{
		{
			UUID:        uuid.NewV4(),
			UserName:    "admin",
			Password:    utils.BcryptHash(initPassword),
			NickName:    "超级管理员",
			Avatar:      "https://img.tuxiangyan.com/zb_users/upload/2023/02/202302091675904134770942.jpg",
			Authority:   system.SysRoleAuthority{AuthorityID: 200, AuthorityName: "系统管理员", ParentID: &superAdmin, DefaultRouter: "dashboard"},
			AuthorityID: 200,
			Phone:       "17700000000",
			Email:       "88888888@qq.com",
			IsActive:    true,
			IsStaff:     true,
		},
	}

	if err := tx.Create(&entities).Error; err != nil {
		return ctx, fmt.Errorf("%s数据初始化失败", system.SysUser{}.TableName())
	}

	authorityEntities, ok := ctx.Value(ctxKey("auth")).([]system.SysRoleAuthority)
	if !ok {
		return ctx, errors.New("创建 [用户-权限] 关联失败, 未找到权限表初始化数据")
	}
	if err := tx.Model(&entities[0]).Association("Authorities").Replace(authorityEntities); err != nil {
		return ctx, err
	}
	return ctx, nil
}
