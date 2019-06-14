package persistence

import (
	"fmt"
	"github.com/docker/distribution/uuid"
	"github.com/splisson/devopstic/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testUser = entities.User{
		Username:  "test",
		Email:     "user@test.com",
		FirstName: "testFirst",
		LastName:  "testLast",
	}
)

func TestCreateGetUser(t *testing.T) {
	username := fmt.Sprintf("testGet%s", uuid.Generate().String())

	t.Run("should create and get user from db", func(t *testing.T) {
		newUser := testUser
		newUser.Username = username
		user, err := testUserStore.CreateUser(newUser)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, user.ID, "id should not be nil")
		assert.NotEmpty(t, user.ID, "id should not be empty")
		user, err = testUserStore.GetUserByUsername(user.Username)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, user, "user exists")
		assert.Equal(t, newUser.Username, user.Username, "same username")
	})
	t.Run("should not create user with same username", func(t *testing.T) {
		uniqueId := uuid.Generate().String()
		newUser := testUser
		newUser.Username = fmt.Sprintf("testGet%s", uniqueId)
		newUser.Email = fmt.Sprintf("testGet%s@test.com", uniqueId)
		_, err := testUserStore.CreateUser(newUser)
		assert.Nil(t, err, "no error expected")
		newUser.Email = fmt.Sprintf("test2Get%s@test.com", uniqueId)
		user, err := testUserStore.CreateUser(newUser)
		assert.NotNil(t, err, "error expected")
		assert.Nil(t, user, "user should be nil")
	})

}
