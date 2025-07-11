package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testRedis *Redis

func setup() {
	envType := "dev"
	testRedis = NewRedis(&envType)
}

func teardown() {
	testRedis.Client.FlushAll(context.Background())
}

func TestNewRedis(t *testing.T) {
	setup()
	defer teardown()

	assert.NotNil(t, testRedis)
	assert.NotNil(t, testRedis.Client)
}

func TestPing(t *testing.T) {
	setup()
	defer teardown()

	response := testRedis.Ping()
	assert.Contains(t, response, "PONG1")
}

func TestStoreEmailHashAndGetEmailFromHash(t *testing.T) {
	setup()
	defer teardown()

	email := "test@example.com"
	hash, err := testRedis.StoreEmailHash(email)
	assert.NoError(t, err)
	assert.NotEmpty(t, hash)

	retrievedEmail, err := testRedis.GetEmailFromHash(hash)
	assert.NoError(t, err)
	assert.Equal(t, email, retrievedEmail)

	// Ensure the key is deleted
	_, err = testRedis.Client.Get(context.Background(), hash).Result()
	assert.Error(t, err)
}

func TestGenerateToken(t *testing.T) {
	setup()
	defer teardown()

	email := "test@example.com"
	token, err := testRedis.GenerateToken(email)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Verify token exists in Redis
	storedEmail, err := testRedis.Client.Get(context.Background(), token).Result()
	assert.NoError(t, err)
	assert.Equal(t, email, storedEmail)
}

func TestCheckSession(t *testing.T) {
	setup()
	defer teardown()

	email := "test@example.com"
	token := "test-token"
	err := testRedis.StoreSession(email, token)
	assert.NoError(t, err)

	exists, err := testRedis.CheckSession(email)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestStoreSession(t *testing.T) {
	setup()
	defer teardown()

	email := "test@example.com"
	token := "test-token"
	err := testRedis.StoreSession(email, token)
	assert.NoError(t, err)

	// Verify session exists
	storedToken, err := testRedis.Client.Get(context.Background(), "session:"+email).Result()
	assert.NoError(t, err)
	assert.Equal(t, token, storedToken)
}

func TestDeleteSession(t *testing.T) {
	setup()
	defer teardown()

	email := "test@example.com"
	token := "test-token"
	err := testRedis.StoreSession(email, token)
	assert.NoError(t, err)

	err = testRedis.DeleteSession(email)
	assert.NoError(t, err)

	// Verify session is deleted
	_, err = testRedis.Client.Get(context.Background(), "session:"+email).Result()
	assert.Error(t, err)
}
