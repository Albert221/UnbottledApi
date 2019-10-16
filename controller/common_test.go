package controller

import (
	"context"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

var (
	johnDoe = &entity.User{
		Base:     entity.Base{ID: uuid.MustParse("bd65f5b2-8563-40c8-8ce6-4d19164fb045")},
		Username: "john.doe",
		Email:    "john.doe@example.com",
		Password: "$2y$12$rf22cbpj6wHhNJf476Wwkee04UrNSv4ZqjwveBChu/cRo1GQkg1s.",
	}
)

func TestGetUser(t *testing.T) {
	t.Run("returns nil when no user in context", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest("GET", "/", nil)
		user := getUser(r)

		assert.Nil(t, user)
	})

	t.Run("returns user when it is present in context", func(t *testing.T) {
		t.Parallel()

		r := httptest.NewRequest("GET", "/", nil)
		r = r.WithContext(context.WithValue(r.Context(), userContextKey{}, johnDoe))

		user := getUser(r)

		assert.Equal(t, johnDoe, user)
	})
}