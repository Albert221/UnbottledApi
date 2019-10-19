package mock

import (
	"github.com/Albert221/UnbottledApi/entity"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) ByID(id uuid.UUID) *entity.User {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.User)
}

func (m *UserRepositoryMock) ByUsernameOrEmail(value string) *entity.User {
	args := m.Called(value)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*entity.User)
}

func (m *UserRepositoryMock) Save(user *entity.User) error {
	args := m.Called(user)

	return args.Error(0)
}
