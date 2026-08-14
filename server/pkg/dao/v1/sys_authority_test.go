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

func setupAuthorityTestContainer(t *testing.T) (context.Context, *gorm.DB, func()) {
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
		&system.SysRoleAuthority{},
		&system.SysUser{},
		&system.SysMenu{},
		&system.SysUserAuthorityRelation{},
	)
	assert.NoError(t, err)

	cleanup := func() {
		mysqlC.Terminate(ctx)
	}

	return ctx, db, cleanup
}

func TestAuthorityDao_AddAuthority(t *testing.T) {
	ctx, db, cleanup := setupAuthorityTestContainer(t)
	defer cleanup()

	authorityDao := NewAuthorityDao(db)

	t.Run("add new authority", func(t *testing.T) {
		authority := system.SysRoleAuthority{
			AuthorityID:   1,
			AuthorityName: "Test Authority",
			ParentID:      func() *uint { u := uint(0); return &u }(),
		}

		apis, err := authorityDao.AddAuthority(ctx, authority)
		assert.NoError(t, err)
		assert.NotNil(t, apis)
	})

	t.Run("add existing authority", func(t *testing.T) {
		authority := system.SysRoleAuthority{
			AuthorityID:   1,
			AuthorityName: "Duplicate Authority",
			ParentID:      func() *uint { u := uint(0); return &u }(),
		}

		_, err := authorityDao.AddAuthority(ctx, authority)
		assert.Error(t, err)
	})
}

func TestAuthorityDao_DeleteAuthority(t *testing.T) {
	ctx, db, cleanup := setupAuthorityTestContainer(t)
	defer cleanup()

	authorityDao := NewAuthorityDao(db)

	parentID := uint(0)
	authority := system.SysRoleAuthority{
		AuthorityID:   2,
		AuthorityName: "Delete Authority",
		ParentID:      &parentID,
	}
	err := db.Create(&authority).Error
	assert.NoError(t, err)

	t.Run("delete existing authority", func(t *testing.T) {
		err := authorityDao.DeleteAuthority(ctx, authority)
		assert.NoError(t, err)
	})

	t.Run("delete non-existent authority", func(t *testing.T) {
		nonExistentAuth := system.SysRoleAuthority{
			AuthorityID: 99999,
		}
		err := authorityDao.DeleteAuthority(ctx, nonExistentAuth)
		assert.Error(t, err)
	})
}

func TestAuthorityDao_GetAuthorityList(t *testing.T) {
	ctx, db, cleanup := setupAuthorityTestContainer(t)
	defer cleanup()

	authorityDao := NewAuthorityDao(db)

	parentID := uint(0)
	authority := system.SysRoleAuthority{
		AuthorityID:   3,
		AuthorityName: "Get Authority",
		ParentID:      &parentID,
	}
	err := db.Create(&authority).Error
	assert.NoError(t, err)

	t.Run("get authority list", func(t *testing.T) {
		list, err := authorityDao.GetAuthorityList(ctx, authority.AuthorityID)
		assert.NoError(t, err)
		assert.NotEmpty(t, list)
	})

	t.Run("get non-existent authority list", func(t *testing.T) {
		_, err := authorityDao.GetAuthorityList(ctx, 99999)
		assert.Error(t, err)
	})
}

func TestAuthorityDao_UpdateAuthority(t *testing.T) {
	ctx, db, cleanup := setupAuthorityTestContainer(t)
	defer cleanup()

	authorityDao := NewAuthorityDao(db)

	parentID := uint(0)
	authority := system.SysRoleAuthority{
		AuthorityID:   4,
		AuthorityName: "Update Authority",
		ParentID:      &parentID,
	}
	err := db.Create(&authority).Error
	assert.NoError(t, err)

	t.Run("update authority", func(t *testing.T) {
		authority.AuthorityName = "Updated Authority"
		updated, err := authorityDao.UpdateAuthority(ctx, authority)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Authority", updated.AuthorityName)
	})

	t.Run("update non-existent authority", func(t *testing.T) {
		nonExistentAuth := system.SysRoleAuthority{
			AuthorityID:   99999,
			AuthorityName: "Non Existent",
		}
		_, err := authorityDao.UpdateAuthority(ctx, nonExistentAuth)
		assert.Error(t, err)
	})
}

func TestAuthorityDao_GetParentAuthorityID(t *testing.T) {
	ctx, db, cleanup := setupAuthorityTestContainer(t)
	defer cleanup()

	authorityDao := NewAuthorityDao(db)

	parentID := uint(0)
	authority := system.SysRoleAuthority{
		AuthorityID:   5,
		AuthorityName: "Parent Authority",
		ParentID:      &parentID,
	}
	err := db.Create(&authority).Error
	assert.NoError(t, err)

	t.Run("get parent authority id", func(t *testing.T) {
		parentID, err := authorityDao.GetParentAuthorityID(ctx, authority.AuthorityID)
		assert.NoError(t, err)
		assert.Equal(t, uint(0), parentID)
	})

	t.Run("get parent authority id for non-existent", func(t *testing.T) {
		_, err := authorityDao.GetParentAuthorityID(ctx, 99999)
		assert.Error(t, err)
	})
}
