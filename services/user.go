package services

import (
	"github.com/splisson/devopstic/entities"
	"github.com/splisson/devopstic/persistence"
)

type UserServiceInterface interface {
	GetUserByUsername(username string) (*entities.User, error)
}

type UserService struct {
	userStore persistence.UserStoreInterface
}

func NewUserService(userStore persistence.UserStoreInterface) *UserService {
	service := new(UserService)
	service.userStore = userStore
	return service
}

func (s *UserService) GetUserByUsername(username string) (*entities.User, error) {
	return nil, nil
}
