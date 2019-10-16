package main

import (
	"github.com/Albert221/UnbottledApi/controller"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/mysql"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"time"
)

type application struct {
	port string

	db    *gorm.DB
	users repository.UserRepository

	jwtAlgo jwt.Algorithm
}

func newApplication(dbDsn, port string) (*application, error) {
	db, err := gorm.Open("mysql", dbDsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	return &application{
		port:    port,
		db:      db,
		users:   mysql.NewUserRepository(db),
		jwtAlgo: jwt.NewHS256([]byte("mTdm6czopftZKezaMAS2BWEo91bCVjNF")),
	}, nil
}

func (a *application) Migrate() {
	a.db.AutoMigrate(entity.User{}, entity.Point{})
}

func (a *application) Serve() error {
	authContr := controller.NewAuthController(a.users, a.jwtAlgo)

	r := mux.NewRouter()
	r.Use(mux.MiddlewareFunc(authContr.AuthenticationMiddleware))

	r.HandleFunc("/auth/authenticate", authContr.AuthenticateHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:" + a.port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
