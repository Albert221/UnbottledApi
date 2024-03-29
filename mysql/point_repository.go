package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type PointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) *PointRepository {
	return &PointRepository{db: db}
}

func (PointRepository) result(point *entity.Point) *entity.Point {
	empty := entity.Point{}
	if empty == *point {
		return nil
	}

	return point
}

func (p *PointRepository) ByID(id uuid.UUID) *entity.Point {
	point := new(entity.Point)
	p.db.Preload("Photo").First(point, "id = ?", id.String())

	return p.result(point)
}

func (p *PointRepository) InArea(lat, lng, radius float32) ([]*entity.Point, error) {
	if radius <= 0 {
		return nil, repository.RadiusNegativeOrZeroErr
	}

	var points []*entity.Point
	// Haversine Formula: https://stackoverflow.com/a/29555137/3158312
	p.db.
		Preload("Photo").
		Where("6371 * 2 * ASIN(SQRT(POWER(SIN((? - ABS(latitude)) * PI() / 180 / 2), 2) "+
			"+ COS(? * PI() / 180 ) * COS(ABS(latitude) * PI() / 180) "+
			"* POWER(SIN((? - (longitude)) * PI() / 180 / 2), 2))) <= ?", lat, lat, lng, radius).
		Find(&points)

	return points, nil
}

func (p *PointRepository) ByAuthorID(authorID uuid.UUID) []*entity.Point {
	var points []*entity.Point
	p.db.Preload("Photo").Where("author_id = ?", authorID.String()).Find(&points)

	return points
}

func (p *PointRepository) Save(point *entity.Point) error {
	return p.db.Save(point).Error
}
