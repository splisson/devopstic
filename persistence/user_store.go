package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/opstic/entities"
)

type UserStoreInterface interface {
	GetUserByUsername(username string) (*entities.User, error)
	CreateUser(user entities.User) (*entities.User, error)
}

type UserDBStore struct {
	db *gorm.DB
}

func NewUserDBStore(db *gorm.DB) *UserDBStore {
	store := new(UserDBStore)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *UserDBStore) GetUserByUsername(username string) (*entities.User, error) {
	user := entities.User{}
	db := s.db.Table("users").Select("*").Where("username = ?", username)
	db.Find(&user)
	return &user, nil
}

func (s *UserDBStore) CreateUser(user entities.User) (*entities.User, error) {
	user.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &user, tx.Commit().Error
}