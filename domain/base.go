package domain

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:binary(16);primary_key"`
	CreatedAt time.Time `json:"created_at"`
}

func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	id := uuid.New()

	return scope.SetColumn("ID", id)
}
