package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type PhotoRepository struct {
	db *gorm.DB
}

func NewPhotoRepository(db *gorm.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (PhotoRepository) result(point *entity.Photo) *entity.Photo {
	empty := entity.Photo{}
	if empty == *point {
		return nil
	}

	return point
}

func (p *PhotoRepository) ByID(id uuid.UUID) *entity.Photo {
	photo := new(entity.Photo)
	p.db.First(photo, "id = ?", id.String())

	return p.result(photo)
}

func (p *PhotoRepository) Save(photo *entity.Photo) error {
	return p.db.Save(photo).Error
}
