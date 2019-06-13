package services

import (
	"github.com/splisson/opstic/persistence"
	"github.com/splisson/opstic/entities"
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


