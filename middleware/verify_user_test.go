package middleware_test

import (
	"booking/internal/cache"
	"booking/internal/constant"
	"booking/internal/util"
	"booking/middleware"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestVerifyUser(t *testing.T) {

	viper.Set("jwt.access_token.expires_at", 1)
	viper.Set("jwt.access_token.secret", "secretKey")
	const expiredAccessTokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiYm9va2luZyIsImV4cCI6MTcxMzUxNzI4MCwibmJmIjoxNzEzNTEzNjgwLCJpYXQiOjE3MTM1MTM2ODAsImp0aSI6IjkxMjE0MTk1LTU0YjktNGY4Yi04NmRmLTRmNTViMjE1ZGJlMyJ9.SS_9HwZoU5VMtZKdwIu6sZJpYXmcIl7UYK1YmOHTQwk"
	app := fiber.New()
	mockRedis := memory.New()
	cache := cache.New(mockRedis)
	util := util.New()
	verifyUserMiddleware := middleware.VerifyUser(cache, util)
	const URL = "/test"
	const username = "example@gmail.com"
	app.Get(URL, verifyUserMiddleware, func(c *fiber.Ctx) error {
		rawAccessToken := c.Locals(constant.LOCAL_KEY_ACCESS_TOKEN)
		if rawAccessToken == nil {
			// status code 500 for debug. when a request rejected by the middleware, should response status code 401 or 403
			return c.Status(fiber.StatusInternalServerError).Send([]byte{})
		}
		return c.Status(fiber.StatusOK).JSON(rawAccessToken)
	})

	t.Run("When a token attached to Authorization header and it has correct format", func(t *testing.T) {
		accessTokenString, _ := util.JWTCreateAccessTokenString(username)
		accessTokenKey := cache.GetAccessTokenKey(username)
		mockRedis.Set(accessTokenKey, []byte(accessTokenString), time.Duration(1*time.Minute))
		defer mockRedis.Delete(accessTokenKey)

		req := httptest.NewRequest(fiber.MethodGet, URL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessTokenString)

		res, _ := app.Test(req)
		assert.NotNil(t, res)
		if res != nil {
			assert.Equal(t, fiber.StatusOK, res.StatusCode)
		}
	})

	t.Run("When a token attached to Authorization header but the token have invalid format", func(t *testing.T) {
		accessTokenString := "invalid_token_format"
		accessTokenKey := cache.GetAccessTokenKey(username)
		mockRedis.Set(accessTokenKey, []byte(accessTokenString), time.Duration(1*time.Minute))
		defer mockRedis.Delete(accessTokenKey)

		req := httptest.NewRequest(fiber.MethodGet, URL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessTokenString)

		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("When a token attached to Authorization header and it has the correct format but the token had expired, ", func(t *testing.T) {
		req := httptest.NewRequest(fiber.MethodGet, URL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+expiredAccessTokenString)

		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("When a token attached to Authorization header and it has correct format but cannot get the token from cache", func(t *testing.T) {
		accessTokenString, _ := util.JWTCreateAccessTokenString(username)

		req := httptest.NewRequest(fiber.MethodGet, URL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessTokenString)

		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
	})

	t.Run("When a token attached to Authorization header and it has correct format but the token in the cache is not the same as header", func(t *testing.T) {
		accessTokenString, _ := util.JWTCreateAccessTokenString(username)
		accessTokenString2 := "dummy_token"
		accessTokenKey := cache.GetAccessTokenKey(username)
		mockRedis.Set(accessTokenKey, []byte(accessTokenString2), time.Duration(1*time.Minute))
		defer mockRedis.Delete(accessTokenKey)

		req := httptest.NewRequest(fiber.MethodGet, URL, nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessTokenString)

		res, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)

	})
}
