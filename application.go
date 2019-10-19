package main

import (
	"fmt"
	"github.com/Albert221/UnbottledApi/controller"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/mysql"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"os"
	"time"
)

type application struct {
	port string

	db     *gorm.DB
	users  repository.UserRepository
	points repository.PointRepository
	photos repository.PhotoRepository

	jwtAlgo jwt.Algorithm
}

func newApplication(dbDsn, port string, secret []byte) (*application, error) {
	db, err := gorm.Open("mysql", dbDsn+"?parseTime=true")
	if err != nil {
		return nil, err
	}

	if err := prepareFilesystem(); err != nil {
		return nil, err
	}

	return &application{
		port:    port,
		db:      db,
		users:   mysql.NewUserRepository(db),
		points:  mysql.NewPointRepository(db),
		photos:  mysql.NewPhotoRepository(db),
		jwtAlgo: jwt.NewHS256(secret),
	}, nil
}

func prepareFilesystem() error {
	return os.MkdirAll("uploads", os.ModePerm)
}

func (a *application) Migrate() {
	a.db.AutoMigrate(entity.User{}, entity.Point{}, entity.Photo{})
}

func (a *application) Serve() error {
	// todo(Albert221): run a cleanup of photos every x hours to remove not used photos

	authContr := controller.NewAuthController(a.users, a.jwtAlgo)
	userContr := controller.NewUserController(a.users)
	pointContr := controller.NewPointController(a.points, a.photos)

	r := mux.NewRouter()
	r.Use(mux.MiddlewareFunc(authContr.AuthenticationMiddleware))

	r.HandleFunc("/auth/authenticate", authContr.AuthenticateHandler).Methods("POST")

	r.HandleFunc("/user", userContr.CreateHandler).Methods("POST")

	r.HandleFunc("/point/{lat},{lng},{radius}", pointContr.GetPointsHandler).Methods("GET")
	r.HandleFunc("/point/photo", pointContr.UploadPhoto).Methods("POST")
	r.HandleFunc("/point", pointContr.AddHandler).Methods("POST")

	addr := "127.0.0.1:" + a.port
	fmt.Println("listening on " + addr)
	srv := &http.Server{
		Handler:      r,
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	return srv.ListenAndServe()
}
