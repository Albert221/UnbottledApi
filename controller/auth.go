package controller

import (
	"context"
	"encoding/json"
	"github.com/Albert221/UnbottledApi/repository"
	valid "github.com/asaskevich/govalidator"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

type AuthController struct {
	users   repository.UserRepository
	jwtAlgo jwt.Algorithm
}

func NewAuthController(users repository.UserRepository, jwtAlgo jwt.Algorithm) *AuthController {
	return &AuthController{users: users, jwtAlgo: jwtAlgo}
}

type jwtPayload struct {
	jwt.Payload
	UserID uuid.UUID `json:"user_id,omitempty"`
}

func (a *AuthController) AuthenticateHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		EmailOrUsername string `json:"email_or_username" valid:"required"`
		Password        string `json:"password" valid:"required"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Request body must be a valid json"})
		return
	}

	if _, err := valid.ValidateStruct(body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	user := a.users.ByUsernameOrEmail(body.EmailOrUsername)
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Bad credentials"})
		return
	}

	payload := jwtPayload{
		UserID: user.ID,
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(time.Now().Add(2 * time.Hour)),
			IssuedAt:       jwt.NumericDate(time.Now()),
		},
	}

	// todo(Albert221): return refresh token too
	accessToken, err := jwt.Sign(payload, a.jwtAlgo)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"access_token": string(accessToken)})
}

type userContextKey struct{}

func (a *AuthController) AuthenticationMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			h.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) < 2 || strings.ToLower(parts[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Authorization header is invalid"})
			return
		}

		jwtToken := parts[1]

		var payload jwtPayload
		payloadValidator := jwt.ValidatePayload(&payload.Payload, jwt.ExpirationTimeValidator(time.Now()))
		_, err := jwt.Verify([]byte(jwtToken), a.jwtAlgo, &payload, payloadValidator)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Given token is invalid or expired"})
			return
		}

		user := a.users.ById(payload.UserID)
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "User for given token does not exist"})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}