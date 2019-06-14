package persistence

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateGetCommit(t *testing.T) {

	t.Run("should create and get commit in db", func(t *testing.T) {
		newCommit := testCommit
		commit, err := testCommitStore.CreateCommit(newCommit)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, commit.ID, "id should not be nil")
		assert.NotEmpty(t, commit.ID, "id should not be empty")
		commits, err := testCommitStore.GetCommits()
		assert.Nil(t, err, "no error")
		assert.NotNil(t, commits, "commits exists")
		assert.NotEmpty(t, commits, "list not empty")
	})
}
