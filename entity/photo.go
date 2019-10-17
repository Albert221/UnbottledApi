package entity

import "github.com/google/uuid"

type Photo struct {
	Base
	AuthorID uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`
	Author   User      `json:"-" gorm:"foreignkey:AuthorID"`
	FileName string    `json:"filename"`
}
