package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"server/pkg/model/system"
)

func setupMenuTestContainer(t *testing.T) (context.Context, *gorm.DB, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "password",
			"MYSQL_DATABASE":      "test_db",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server - GPL"),
	}

	mysqlC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)

	ip, err := mysqlC.Host(ctx)
	assert.NoError(t, err)

	port, err := mysqlC.MappedPort(ctx, "3306")
	assert.NoError(t, err)

	dsn := "root:password@tcp(" + ip + ":" + port.Port() + ")/test_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&system.SysMenu{},
		&system.SysRoleAuthority{},
		&system.SysMenuAuthorityRelation{},
	)
	assert.NoError(t, err)

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}

	return ctx, db, cleanup
}

func TestMenuDao_AddMenu(t *testing.T) {
	ctx, db, cleanup := setupMenuTestContainer(t)
	defer cleanup()

	menuDao := NewMenuDao(db)

	t.Run("add new menu", func(t *testing.T) {
		menu := system.SysMenu{
			Name:      "Test Menu",
			Path:      "/test",
			ParentID:  0,
			Sort:      1,
			Component: "test",
		}

		err := menuDao.AddMenu(ctx, menu)
		assert.NoError(t, err)
	})

	t.Run("add existing menu", func(t *testing.T) {
		menu := system.SysMenu{
			Name:      "Test Menu",
			Path:      "/test2",
			ParentID:  0,
			Sort:      2,
			Component: "test2",
		}

		err := menuDao.AddMenu(ctx, menu)
		assert.Error(t, err)
	})

	t.Run("add menu with non-existent parent", func(t *testing.T) {
		menu := system.SysMenu{
			Name:      "Child Menu",
			Path:      "/child",
			ParentID:  99999,
			Sort:      3,
			Component: "child",
		}

		err := menuDao.AddMenu(ctx, menu)
		assert.Error(t, err)
	})
}

func TestMenuDao_UpdateMenu(t *testing.T) {
	ctx, db, cleanup := setupMenuTestContainer(t)
	defer cleanup()

	menuDao := NewMenuDao(db)

	menu := system.SysMenu{
		Name:      "Update Menu",
		Path:      "/update",
		ParentID:  0,
		Sort:      1,
		Component: "update",
	}
	err := db.Create(&menu).Error
	assert.NoError(t, err)

	t.Run("update menu", func(t *testing.T) {
		menu.Name = "Updated Menu"
		err := menuDao.UpdateMenu(ctx, menu)
		assert.NoError(t, err)

		var updatedMenu system.SysMenu
		err = db.First(&updatedMenu, menu.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Updated Menu", updatedMenu.Name)
	})
}

func TestMenuDao_DeleteMenu(t *testing.T) {
	ctx, db, cleanup := setupMenuTestContainer(t)
	defer cleanup()

	menuDao := NewMenuDao(db)

	menu := system.SysMenu{
		Name:      "Delete Menu",
		Path:      "/delete",
		ParentID:  0,
		Sort:      1,
		Component: "delete",
	}
	err := db.Create(&menu).Error
	assert.NoError(t, err)

	t.Run("delete menu", func(t *testing.T) {
		err := menuDao.DeleteMenu(ctx, menu.ID)
		assert.NoError(t, err)
	})

	t.Run("delete non-existent menu", func(t *testing.T) {
		err := menuDao.DeleteMenu(ctx, 99999)
		assert.Error(t, err)
	})

	t.Run("delete menu with children", func(t *testing.T) {
		parentMenu := system.SysMenu{
			Name:      "Parent Menu",
			Path:      "/parent",
			ParentID:  0,
			Sort:      2,
			Component: "parent",
		}
		err := db.Create(&parentMenu).Error
		assert.NoError(t, err)

		childMenu := system.SysMenu{
			Name:      "Child Menu",
			Path:      "/child",
			ParentID:  parentMenu.ID,
			Sort:      3,
			Component: "child",
		}
		err = db.Create(&childMenu).Error
		assert.NoError(t, err)

		err = menuDao.DeleteMenu(ctx, parentMenu.ID)
		assert.Error(t, err)
	})
}

func TestMenuDao_GetMenuAuthority(t *testing.T) {
	ctx, db, cleanup := setupMenuTestContainer(t)
	defer cleanup()

	menuDao := NewMenuDao(db)

	menu := system.SysMenu{
		Name:      "Authority Menu",
		Path:      "/authority",
		ParentID:  0,
		Sort:      1,
		Component: "authority",
	}
	err := db.Create(&menu).Error
	assert.NoError(t, err)

	authority := system.SysRoleAuthority{
		AuthorityID:   1,
		AuthorityName: "Test Authority",
	}
	err = db.Create(&authority).Error
	assert.NoError(t, err)

	menuAuthority := system.SysMenuAuthorityRelation{
		MenuID:      string(rune('0' + menu.ID)),
		AuthorityID: "1",
	}
	err = db.Create(&menuAuthority).Error
	assert.NoError(t, err)

	t.Run("get menu authority", func(t *testing.T) {
		menus, err := menuDao.GetMenuAuthority(ctx, authority.AuthorityID)
		assert.NoError(t, err)
		assert.NotEmpty(t, menus)
	})

	t.Run("get menu authority for non-existent authority", func(t *testing.T) {
		menus, err := menuDao.GetMenuAuthority(ctx, 99999)
		assert.NoError(t, err)
		assert.Empty(t, menus)
	})
}
