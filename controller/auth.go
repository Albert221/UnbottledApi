package controller

import (
	"encoding/json"
	"github.com/Albert221/UnbottledApi/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	users repository.UserRepository
}

func NewAuthController(users repository.UserRepository) *AuthController {
	return &AuthController{users: users}
}

func (a *AuthController) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		EmailOrUsername string `json:"email_or_username"`
		Password        string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Request body must be a valid json"})
		return
	}

	user := a.users.ByUsernameOrEmail(body.EmailOrUsername)
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Bad credentials"})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Authenticated successfully"})

	// todo: create and return jwt token
}
