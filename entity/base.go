package entity

import (
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"time"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
}

func (b *Base) BeforeCreate(scope *gorm.Scope) error {
	if b.ID.ID() == 0 {
		if err := scope.SetColumn("ID", uuid.New()); err != nil {
			return err
		}
	}

	if b.CreatedAt.IsZero() {
		return scope.SetColumn("CreatedAt", time.Now().UTC())
	}

	return nil
}
