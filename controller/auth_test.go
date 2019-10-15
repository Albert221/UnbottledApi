package controller

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	johnDoe = &entity.User{
		Username: "john.doe",
		Email:    "john.doe@example.com",
		Password: "$2y$12$rf22cbpj6wHhNJf476Wwkee04UrNSv4ZqjwveBChu/cRo1GQkg1s.",
	}
)

func TestAuthController_AuthenticateHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	body := strings.NewReader(`{"email_or_username": "john.doe", "password": "password"}`)
	r := httptest.NewRequest("POST", "/auth/authenticate", body)

	userRepoMock := new(userRepositoryMock)
	userRepoMock.On("ByUsernameOrEmail", "john.doe").Return(johnDoe)
	userRepoMock.On("ByUsernameOrEmail", "john.doe@example.com").Return(johnDoe)
	userRepoMock.On("ByUsernameOrEmail", "idontexist").Return(nil)

	contr := NewAuthController(userRepoMock)

	contr.AuthenticateHandler(rr, r)

	cases := []struct {
		Name  string
		Body  string
		Check func(t *testing.T, r *httptest.ResponseRecorder)
	}{
		{
			Name: "correctly authenticates john by username",
			Body: `{"email_or_username": "john.doe", "password": "password"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.JSONEq(t, `{"message": "Authenticated successfully"}`, r.Body.String())
			},
		},
		{
			Name: "correctly authenticates john by email",
			Body: `{"email_or_username": "john.doe@example.com", "password": "password"}`,
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				assert.JSONEq(t, `{"message": "Authenticated successfully"}`, r.Body.String())
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
			Name: "correctly fails when body isn't a valid json",
			Body: "hey im not a valid json",
			Check: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				assert.JSONEq(t, `{"error": "Request body must be a valid json"}`, r.Body.String())
			},
		},
	}

	for _, aCase := range cases {
		aCase := aCase // capture range variable
		t.Run(aCase.Name, func(t *testing.T) {
			t.Parallel()

			rr := httptest.NewRecorder()
			body := strings.NewReader(aCase.Body)
			r := httptest.NewRequest("POST", "/auth/authenticate", body)

			contr.AuthenticateHandler(rr, r)

			aCase.Check(t, rr)
		})
	}
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
