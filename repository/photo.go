package repository

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
)

type PhotoRepository interface {
	ByID(id uuid.UUID) *entity.Photo
	Save(photo *entity.Photo) error
}
