package repository

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	ByID(id uuid.UUID) *entity.User
	ByUsernameOrEmail(value string) *entity.User
	Save(user *entity.User) error
}
