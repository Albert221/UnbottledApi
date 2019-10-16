package controller

import (
	"encoding/json"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserController struct {
	users repository.UserRepository
}

func NewUserController(users repository.UserRepository) *UserController {
	return &UserController{users: users}
}

func (u *UserController) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username" valid:"required"`
		Email    string `json:"email" valid:"required,email"`
		Password string `json:"password" valid:"required"`
	}

	if err := decodeAndValidateBody(&body, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// todo(Albert221): send activation email

	user := &entity.User{
		Username: body.Username,
		Email:    body.Email,
		Password: string(password),
	}

	if err := u.users.Save(user); err != nil {
		// todo(Albert221): check for unique constraint fail
		//if err.

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{"user": user})
}
