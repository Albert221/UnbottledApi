package entity

import (
	"github.com/google/uuid"
	"time"
)

type Base struct {
	ID        uuid.UUID `json:"id" gorm:"type:char(36);primary_key"`
	CreatedAt time.Time `json:"created_at"`
}

//func (b *Base) BeforeCreate(scope *gorm.Scope) error {
//	id := uuid.New()
//
//	// todo: only do this if it's not already set
//	if err := scope.SetColumn("ID", id); err != nil {
//		return err
//	}
//
//	return scope.SetColumn("CreatedAt", time.Now())
//}
