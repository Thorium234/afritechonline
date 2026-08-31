package mikrotik_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/Thorium234/afritechonline/backend/internal/mikrotik"
)

func TestConnectAndIdentity(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	client := mikrotik.NewClient("127.0.0.1", "admin", "password", 5*time.Second)
	identity, version, err := client.TestConnection(context.Background())
	assert.NoError(t, err)
	assert.NotEmpty(t, identity)
	assert.NotEmpty(t, version)
}

func TestUserLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	client := mikrotik.NewClient("127.0.0.1", "admin", "password", 5*time.Second)
	userService := mikrotik.NewUserService(client)

	err := userService.CreateUser(context.Background(), "testuser", "testpass", "default")
	assert.NoError(t, err)

	err = userService.DisableUser(context.Background(), "testuser")
	assert.NoError(t, err)

	err = userService.EnableUser(context.Background(), "testuser")
	assert.NoError(t, err)

	err = userService.DeleteUser(context.Background(), "testuser")
	assert.NoError(t, err)
}
