package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type RatingRepository struct {
	db *gorm.DB
}

func NewRatingRepository(db *gorm.DB) *RatingRepository {
	return &RatingRepository{db: db}
}

func (r *RatingRepository) ByPointID(pointID uuid.UUID) []*entity.Rating {
	var ratings []*entity.Rating
	r.db.Where("point_id = ?", pointID.String()).Find(&ratings)

	return ratings
}

func (r *RatingRepository) Save(rating *entity.Rating) error {
	return r.db.Save(rating).Error
}
