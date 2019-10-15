package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) ById(id uuid.UUID) *entity.User {
	user := new(entity.User)
	u.db.First(&user, "id = ?", id.String())

	return user
}

func (u *UserRepository) ByUsernameOrEmail(value string) *entity.User {
	user := new(entity.User)
	u.db.First(user, "username = ? OR email = ?", value, value)

	return user
}
