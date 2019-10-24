package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	johnDoe = &entity.User{
		Base: entity.Base{
			ID:        uuid.MustParse("80cc7eb0-c960-43af-9b90-e421ddca52e5"),
			CreatedAt: time.Now(),
		},
		Email:    "john.doe@example.com",
		Username: "john.doe",
		Password: "test",
	}
	mikeDoe = &entity.User{
		Base: entity.Base{
			ID:        uuid.MustParse("76823743-19b1-4efb-b2e8-d1af41f58f33"),
			CreatedAt: time.Now(),
		},
		Email:    "mike.doe@example.com",
		Username: "mike.doe",
		Password: "test123",
	}
)

func setupUsersDb() *gorm.DB {
	db, err := gorm.Open("mysql", os.Getenv("TEST_DB_DSN")+"?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(entity.User{})
	db.Exec("TRUNCATE TABLE `users`")

	return db
}

func TestUserRepository_ById(t *testing.T) {
	db := setupUsersDb()
	defer db.Close()

	db.Create(johnDoe)
	db.Create(mikeDoe)

	repo := NewUserRepository(db)

	t.Run("returns correctly john", func(t *testing.T) {
		user := repo.ByID(johnDoe.ID)
		assert.Equal(t, johnDoe.ID, user.ID)
	})

	t.Run("returns correctly nothing", func(t *testing.T) {
		user := repo.ByID(uuid.MustParse("f923f3e3-94c0-43c3-83e2-f9a772e16f23"))
		assert.Nil(t, user)
	})
}

func TestUserRepository_ByUsernameOrEmail(t *testing.T) {
	db := setupUsersDb()
	defer db.Close()

	db.Create(johnDoe)
	db.Create(mikeDoe)

	repo := NewUserRepository(db)

	t.Run("returns correctly john by username", func(t *testing.T) {
		user := repo.ByUsernameOrEmail("john.doe")
		assert.Equal(t, johnDoe.ID, user.ID)
	})

	t.Run("returns correctly john by email", func(t *testing.T) {
		user := repo.ByUsernameOrEmail("john.doe@example.com")
		assert.Equal(t, johnDoe.ID, user.ID)
	})

	t.Run("returns correctly nothing", func(t *testing.T) {
		user := repo.ByUsernameOrEmail("invalid-value")
		assert.Nil(t, user)
	})
}
