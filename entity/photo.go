package entity

import (
	"github.com/Albert221/UnbottledApi/storage"
	"github.com/google/uuid"
	"path"
)

type Photo struct {
	Base
	AuthorID uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`
	Author   User      `json:"-" gorm:"foreignkey:AuthorID"`
	FileName string    `json:"-"`
	Url      string    `json:"url" gorm:"-"`
}

func (p *Photo) PopulateUrl(host string) {
	if p.FileName == "" {
		return
	}

	p.Url = "http://" + path.Join(host, storage.UploadsPath, p.FileName)
}
