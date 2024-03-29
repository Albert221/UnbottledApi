package entity

import "github.com/google/uuid"

type Rating struct {
	Base
	AuthorID uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`
	Author   User      `json:"-" gorm:"foreignkey:AuthorID"`
	PointID  uuid.UUID `json:"point_id" gorm:"type:char(36);not null"`
	Point    Point     `json:"-" gorm:"foreignkey:PointID"`

	Taste uint32 `json:"taste"`
}