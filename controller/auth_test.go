package controller

import (
	"context"
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAuthController_AuthenticateHandler(t *testing.T) {
	userRepoMock := new(userRepositoryMock)
	userRepoMock.On("ByUsernameOrEmail", "john.doe").Return(johnDoe)
	userRepoMock.On("ByUsernameOrEmail", "john.doe@example.com").Return(johnDoe)
	userRepoMock.On("ByUsernameOrEmail", "idontexist").Return(nil)

	contr := NewAuthController(userRepoMock, jwt.None())

	tests := []struct {
		Name  string
		Body  string
		Check func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			Name: "correctly authenticates john by username",
			Body: `{"email_or_username": "john.doe", "password": "password"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Contains(t, r.Body.String(), "access_token")
			},
		},
		{
			Name: "correctly authenticates john by email",
			Body: `{"email_or_username": "john.doe@example.com", "password": "password"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.Contains(t, r.Body.String(), "access_token")
			},
		},
		{
			Name: "correctly fails authenticating john with wrong password",
			Body: `{"email_or_username": "john.doe", "password": "wrong password"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				assert.JSONEq(t, `{"error": "Bad credentials"}`, r.Body.String())
			},
		},
		{
			Name: "correctly fails authenticating with not existing user",
			Body: `{"email_or_username": "idontexist", "password": "testtest"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				assert.JSONEq(t, `{"error": "Bad credentials"}`, r.Body.String())
			},
		},
		{
			Name: "correctly fails when body is not a valid json",
			Body: "hey im not a valid json",
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				assert.JSONEq(t, `{"error": "Request body must be a valid json"}`, r.Body.String())
			},
		},
		{
			Name: "correctly fails when body is not a valid schema",
			Body: `{"some_field": "yes"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
			},
		},
	}

	for _, test := range tests {
		test := test // capture range variable
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			body := strings.NewReader(test.Body)
			r := httptest.NewRequest("POST", "/auth/authenticate", body)

			contr.AuthenticateHandler(rr, r)

			test.Check(t, rr)
		})
	}
}

func TestAuthController_AuthenticationMiddleware(t *testing.T) {
	okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// do nothing, it returns HTTP 200 and if there is user in context - prints its id
		user := r.Context().Value(userContextKey{})
		if user != nil {
			_, _ = io.WriteString(w, user.(*entity.User).ID.String())
		}
	})

	userRepoMock := new(userRepositoryMock)
	userRepoMock.On("ById", uuid.MustParse("5eb2dd69-a43c-416f-a8ca-90eeb15c12e7")).
		Return(nil)
	userRepoMock.On("ById", johnDoe.ID).Return(johnDoe)

	contr := NewAuthController(userRepoMock, jwt.None())

	tests := []struct {
		Name       string
		AuthHeader string
		Check      func(t *testing.T, rr *httptest.ResponseRecorder)
	}{
		{
			Name: "returns ok when no authorization",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
			},
		},
		{
			Name:       "returns error when authorization not bearer",
			AuthHeader: "Basic 123",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			Name:       "returns error when authorization bearer but empty",
			AuthHeader: "Bearer ",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			Name:       "returns error when token is invalid",
			AuthHeader: "Bearer invalidtoken",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			Name: "returns error when token is expired",
			AuthHeader: "Bearer ew0KICAiYWxnIjogIm5vbmUiLA0KICAidHlwIjogIkpXVCINCn0" +
				".eyJpYXQiOjE1MTYyMzkwMjIsImV4cCI6MTU3MTE3NjI1N30.",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
			},
		},
		{
			Name: "returns error when token user does not exist",
			// token with user_id: 5eb2dd69-a43c-416f-a8ca-90eeb15c12e7
			AuthHeader: "Bearer ew0KICAiYWxnIjogIm5vbmUiLA0KICAidHlwIjogIkpXVCINCn0.eyJpYXQiOjE1MTYyMzkwMj" +
				"IsImV4cCI6NDEwMjQ0ODQ2MSwidXNlcl9pZCI6IjVlYjJkZDY5LWE0M2MtNDE2Zi1hOGNhLTkwZWViMTVjMTJlNyJ9.",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusUnauthorized, rr.Code)
				assert.JSONEq(t, `{"error": "User for given token does not exist"}`, rr.Body.String())
			},
		},
		{
			Name: "returns ok when token is valid and user exists",
			// token with user_id: bd65f5b2-8563-40c8-8ce6-4d19164fb045
			AuthHeader: "Bearer ew0KICAiYWxnIjogIm5vbmUiLA0KICAidHlwIjogIkpXVCINCn0.eyJpYXQiOjE1MTYyMzkwMj" +
				"IsImV4cCI6NDEwMjQ0ODQ2MSwidXNlcl9pZCI6ImJkNjVmNWIyLTg1NjMtNDBjOC04Y2U2LTRkMTkxNjRmYjA0NSJ9.",
			Check: func(t *testing.T, rr *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, rr.Code)
				assert.Equal(t, johnDoe.ID.String(), rr.Body.String())
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			r := httptest.NewRequest("GET", "/", nil)
			if test.AuthHeader != "" {
				r.Header.Set("Authorization", test.AuthHeader)
			}

			middleware := contr.AuthenticationMiddleware(okHandler)
			middleware.ServeHTTP(rr, r)

			test.Check(t, rr)
		})
	}
}

func TestAuthController_getUser(t *testing.T) {
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

type userRepositoryMock struct {
	mock.Mock
}

func (m *userRepositoryMock) ById(id uuid.UUID) *entity.User {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.User)
}

func (m *userRepositoryMock) ByUsernameOrEmail(value string) *entity.User {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.User)
}
