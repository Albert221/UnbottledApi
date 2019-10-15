package mysql

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	gdansk = &entity.Point{
		Base: entity.Base{
			ID:        uuid.MustParse("fa8a3e1d-05bd-42d7-a8c7-3cffe9471b88"),
			CreatedAt: time.Now(),
		},
		Latitude:  52.229788,
		Longitude: 21.011729,
	}
	warsaw = &entity.Point{
		Base: entity.Base{
			ID:        uuid.MustParse("041bf81f-a171-49a1-9da4-bafa1bf0e8ca"),
			CreatedAt: time.Now(),
		},
		Latitude:  54.351881,
		Longitude: 18.646265,
	}
)

func setupPointsDb() *gorm.DB {
	db, err := gorm.Open("mysql", os.Getenv("TEST_DB_DSN")+"?parseTime=true")
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(entity.Point{})
	db.Exec("TRUNCATE TABLE `points`")

	return db
}

func TestPointRepository_ById(t *testing.T) {
	db := setupPointsDb()
	defer db.Close()

	db.Create(gdansk)
	db.Create(warsaw)

	repo := NewPointRepository(db)

	t.Run("returns correctly warsaw", func(t *testing.T) {
		point := repo.ById(warsaw.ID)

		assert.NotNil(t, point)
		assert.Equal(t, warsaw.ID, point.ID)
	})

	t.Run("returns correctly nothing", func(t *testing.T) {
		point := repo.ById(uuid.MustParse("28c80525-2617-4e2c-9aba-a7e1caff1347")) // random uuid

		assert.Zero(t, *point)
	})
}

func TestPointRepository_InArea(t *testing.T) {
	db := setupPointsDb()
	defer db.Close()

	db.Create(gdansk)
	db.Create(warsaw)

	repo := NewPointRepository(db)

	t.Run("returns correctly warsaw", func(t *testing.T) {
		points, err := repo.InArea(warsaw.Latitude, warsaw.Longitude, 100)

		assert.NoError(t, err)
		assert.Len(t, points, 1)
		assert.Equal(t, warsaw.ID, points[0].ID)
	})

	t.Run("returns correctly both warsaw and gdansk", func(t *testing.T) {
		points, err := repo.InArea(warsaw.Latitude, warsaw.Longitude, 1000)

		assert.NoError(t, err)
		assert.Len(t, points, 2)
	})

	t.Run("returns nothing", func(t *testing.T) {
		points, err := repo.InArea(0, 0, 10)

		assert.NoError(t, err)
		assert.Len(t, points, 0)
	})

	t.Run("returns error on negative radius", func(t *testing.T) {
		_, err := repo.InArea(0, 0, -15) // (0, 0) is faaaaaar from Gdansk and Warsaw

		assert.Error(t, err, repository.RadiusNegativeOrZeroErr.Error())
	})
}
