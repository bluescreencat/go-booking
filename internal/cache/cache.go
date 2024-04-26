package cache

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

type cache struct {
	redis fiber.Storage
}

type ICache interface {
	SaveRefreshTokenString(username string, refreshTokenString string) error
	SaveAccessTokenString(username string, accessTokenString string) error
	GetRefreshTokenString(username string) (string, error)
	GetRefreshTokenKey(username string) string
	GetAccessTokenString(username string) (string, error)
	GetAccessTokenKey(username string) string
}

func New(redis fiber.Storage) *cache {
	return &cache{redis: redis}
}

func (cache *cache) GetRefreshTokenKey(username string) string {
	return viper.GetString("redis.prefix_key.refresh_token") + ":" + username
}

func (cache *cache) GetAccessTokenKey(username string) string {
	return viper.GetString("redis.prefix_key.access_token") + ":" + username
}

func (cache *cache) SaveRefreshTokenString(username string, refreshTokenString string) error {
	return cache.redis.Set(
		cache.GetRefreshTokenKey(username),
		[]byte(refreshTokenString),
		viper.GetDuration("jwt.refresh_token.duration"),
	)
}

func (cache *cache) SaveAccessTokenString(username string, accessTokenString string) error {
	return cache.redis.Set(
		cache.GetAccessTokenKey(username),
		[]byte(accessTokenString),
		viper.GetDuration("jwt.access_token.duration"),
	)
}

func (cache *cache) GetRefreshTokenString(username string) (string, error) {
	refreshTokenBytes, err := cache.redis.Get(cache.GetRefreshTokenKey(username))
	if err != nil {
		return "", err
	}

	return string(refreshTokenBytes), nil
}

func (cache *cache) GetAccessTokenString(username string) (string, error) {
	accessTokenBytes, err := cache.redis.Get(cache.GetAccessTokenKey(username))
	if err != nil {
		return "", err
	}

	return string(accessTokenBytes), nil
}
