package controller

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/Albert221/UnbottledApi/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
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
		writeJSON(w, map[string]string{"error": err.Error()})
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
		Active:   true,
	}

	if err := u.users.Save(user); err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			w.WriteHeader(http.StatusBadRequest)
			if strings.Contains(err.Error(), "_username") {
				writeJSON(w, map[string]string{"error": "That username is already taken"})
			} else { // email
				writeJSON(w, map[string]string{"error": "That email is already taken"})
			}
			return
		}

		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	writeJSON(w, map[string]interface{}{"user": user})
}
