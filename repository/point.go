package repository

import (
	"errors"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
)

var (
	RadiusNegativeOrZeroErr = errors.New("radius must be bigger than zero")
)

type PointRepository interface {
	ByID(id uuid.UUID) *entity.Point
	InArea(lat, lng, radius float32) ([]*entity.Point, error)
	ByAuthorID(authorID uuid.UUID) []*entity.Point
	Save(point *entity.Point) error
}
