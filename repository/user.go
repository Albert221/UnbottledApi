package repository

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
)

type UserRepository interface {
	ById(id uuid.UUID) *entity.User
	ByUsernameOrEmail(value string) *entity.User
}
