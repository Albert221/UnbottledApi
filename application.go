package main

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"time"
)

type application struct {
	db *gorm.DB
}

func newApplication(dbDsn string) (*application, error) {
	db, err := gorm.Open("mysql", dbDsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	return &application{
		db: db,
	}, nil
}

func (a *application) Migrate() {
	a.db.AutoMigrate(entity.User{}, entity.Point{})
}

func (a *application) Serve() error {
	r := mux.NewRouter()

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
