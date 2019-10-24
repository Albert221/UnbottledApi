package entity

import (
	"github.com/Albert221/UnbottledApi/storage"
	"github.com/google/uuid"
	"net/url"
	"path"
)

type Photo struct {
	Base
	AuthorID uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`
	Author   User      `json:"-" gorm:"foreignkey:AuthorID"`
	FileName string    `json:"-"`
	Url      string    `json:"url" gorm:"-"`
}

func (p *Photo) PopulateUrl(reqUrl *url.URL) {
	if p.FileName == "" {
		return
	}

	photoUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		User:   reqUrl.User,
		Host:   reqUrl.Host,
		Path:   path.Join(storage.UploadsPath, p.FileName),
	}

	p.Url = photoUrl.String()
}
