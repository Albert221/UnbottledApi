package controller

import (
	"context"
	"github.com/Albert221/UnbottledApi/repository"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
	"time"
)

// todo(Albert221): DRY here!!!

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

	if err := decodeAndValidateBody(&body, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	user := a.users.ByUsernameOrEmail(body.EmailOrUsername)
	if user == nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "Bad credentials"})
		return
	}

	if !user.Active {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "User is not active"})
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

	writeJSON(w, map[string]interface{}{
		"access_token": string(accessToken),
		"user":         user,
	})
}

func (a *AuthController) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// todo(Albert221): use refresh token for its refreshment, not the ordinary token
	var body struct {
		OldToken string `json:"old_token" valid:"required"`
	}

	if err := decodeAndValidateBody(&body, r); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	var payload jwtPayload
	_, err := jwt.Verify([]byte(body.OldToken), a.jwtAlgo, &payload)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]string{"error": "Given token is invalid"})
		return
	}

	user := a.users.ByID(payload.UserID)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]string{"error": "User for given token does not exist"})
		return
	}

	if !user.Active {
		w.WriteHeader(http.StatusUnauthorized)
		writeJSON(w, map[string]string{"error": "User is not active"})
		return
	}

	payload = jwtPayload{
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

	writeJSON(w, map[string]interface{}{
		"access_token": string(accessToken),
		"user":         user,
	})
}

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
			writeJSON(w, map[string]string{"error": "Authorization header is invalid"})
			return
		}

		jwtToken := parts[1]

		var payload jwtPayload
		payloadValidator := jwt.ValidatePayload(&payload.Payload, jwt.ExpirationTimeValidator(time.Now()))
		_, err := jwt.Verify([]byte(jwtToken), a.jwtAlgo, &payload, payloadValidator)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, map[string]string{"error": "Given token is invalid or expired"})
			return
		}

		user := a.users.ByID(payload.UserID)
		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, map[string]string{"error": "User for given token does not exist"})
			return
		}

		if !user.Active {
			w.WriteHeader(http.StatusUnauthorized)
			writeJSON(w, map[string]string{"error": "User is not active"})
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}
