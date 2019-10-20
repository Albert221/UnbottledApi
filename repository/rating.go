package repository

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
)

type RatingRepository interface {
	ByPointID(pointID uuid.UUID) []*entity.Rating
	Save(rating *entity.Rating) error
}
