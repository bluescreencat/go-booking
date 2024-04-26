package cache_test

import (
	"booking/internal/cache"
	"testing"
	"time"

	"github.com/gofiber/storage/memory/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const ACCESS_TOKEN_STRING = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiYm9va2luZyIsImV4cCI6MTcxMzUxNzI4MCwibmJmIjoxNzEzNTEzNjgwLCJpYXQiOjE3MTM1MTM2ODAsImp0aSI6IjkxMjE0MTk1LTU0YjktNGY4Yi04NmRmLTRmNTViMjE1ZGJlMyJ9.SS_9HwZoU5VMtZKdwIu6sZJpYXmcIl7UYK1YmOHTQwk"
const USERNAME = "example@gmail.com"
const PREFIX_KEY = "dummy_prefix"

// change it later when structure of the refrsh token is not the same as the access token
const REFRESH_TOKEN_STRING = ACCESS_TOKEN_STRING

func TestSaveRefreshToken(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("jwt.refresh_token.duration", time.Duration.Minutes(1))
	viper.Set("redis.prefix_key.refresh_token", PREFIX_KEY)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	err := cache.SaveRefreshTokenString(USERNAME, REFRESH_TOKEN_STRING)

	assert.Nil(t, err)
}

func TestSaveAccessToken(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("jwt.access_token.duration", time.Duration.Minutes(1))
	viper.Set("redis.prefix_key.access_token", PREFIX_KEY)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	err := cache.SaveAccessTokenString(USERNAME, ACCESS_TOKEN_STRING)

	assert.Nil(t, err)
}

func TestGetRefreshToken(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("jwt.refresh_token.duration", time.Duration.Minutes(1))
	viper.Set("redis.prefix_key.refresh_token", PREFIX_KEY)
	cache.SaveRefreshTokenString(USERNAME, REFRESH_TOKEN_STRING)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	refreshTokenString, err := cache.GetRefreshTokenString(USERNAME)
	assert.Nil(t, err)
	assert.NotEmpty(t, refreshTokenString)
	if err == nil {
		assert.Equal(t, refreshTokenString, REFRESH_TOKEN_STRING)
	}
}

func TestGetRefreshTokenKey(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("redis.prefix_key.refresh_token", PREFIX_KEY)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	refreshTokenKey := cache.GetRefreshTokenKey(USERNAME)

	assert.Equal(t, PREFIX_KEY+":"+USERNAME, refreshTokenKey)
}

func TestGetAccessTokenString(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("jwt.access_token.duration", time.Duration.Minutes(1))
	viper.Set("redis.prefix_key.access_token", PREFIX_KEY)
	const username = "example@gmail.com"
	cache.SaveRefreshTokenString(username, ACCESS_TOKEN_STRING)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	refreshTokenString, err := cache.GetRefreshTokenString(username)
	assert.Nil(t, err)
	assert.NotEmpty(t, refreshTokenString)
	if err == nil {
		assert.Equal(t, refreshTokenString, ACCESS_TOKEN_STRING)
	}
}

func TestGetAccessTokenKey(t *testing.T) {
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	viper.Set("redis.prefix_key.access_token", PREFIX_KEY)
	defer mockRedis.Delete(PREFIX_KEY + ":" + USERNAME)

	accessTokenKey := cache.GetAccessTokenKey(USERNAME)

	assert.Equal(t, PREFIX_KEY+":"+USERNAME, accessTokenKey)
}
