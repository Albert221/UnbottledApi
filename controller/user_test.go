package controller

import (
	rmock "github.com/Albert221/UnbottledApi/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserController_CreateHandler(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/", strings.NewReader(
		`{"username": "john.doe", "email": "john.doe@example.com", "password": "tak123"}`))

	userRepoMock := new(rmock.UserRepositoryMock)
	userRepoMock.On("Save", mock.AnythingOfType("*entity.User")).Return(nil)

	contr := NewUserController(userRepoMock)

	contr.CreateHandler(rr, r)

	assert.NotEmpty(t, rr.Body.String())
}
