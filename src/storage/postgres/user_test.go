package postgres

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userRepo UserRepo

func TestMain(m *testing.M) {
	PrepareFixtures()

	userRepo = UserRepo{GetDB()}

	os.Exit(m.Run())
}

// nolint: paralleltest
func TestProfileRepo_GetProfilesByEmail(t *testing.T) {
	LoadFixtures()

	email := "email1@test.com"

	user, err := userRepo.GetUserByEmail(context.Background(), email)
	if !assert.NoError(t, err) {
		return
	}

	assert.Equal(t, user.Email, email)
}
