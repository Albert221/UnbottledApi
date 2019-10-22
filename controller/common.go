package controller

import (
	"encoding/json"
	"errors"
	"github.com/Albert221/UnbottledApi/entity"
	valid "github.com/asaskevich/govalidator"
	"net/http"
)

type userContextKey struct{}

func getUser(r *http.Request) *entity.User {
	user := r.Context().Value(userContextKey{})
	if user == nil {
		return nil
	}

	return user.(*entity.User)
}

func decodeAndValidateBody(body interface{}, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return errors.New("request body must be a valid json")
	}

	_, err := valid.ValidateStruct(body)

	return err
}

func writeJSON(w http.ResponseWriter, body interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(body)
}
