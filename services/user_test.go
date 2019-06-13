package services

import (
	"github.com/google/uuid"
	"github.com/splisson/opstic/persistence"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var username string = "admin@weyv.com"

func TestGetUserByUsername(t *testing.T) {
	dbId := uuid.New().String()
	db, dbFilepath := persistence.NewSQLiteConnection(dbId)
	defer db.Close()
	userStore := persistence.NewUserDBStore(db)
	userService := NewUserService(userStore)

	t.Run("Returns user", func(t *testing.T) {
		user, err := userService.GetUserByUsername(username)
		assert.Nil(t, err, "No error expected")
		assert.NotNil(t, user, "User should be defined")
	})

	os.Remove(dbFilepath)
}
