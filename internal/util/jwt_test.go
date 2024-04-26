package util_test

import (
	"booking/internal/util"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestJWTCreateRefreshToken(t *testing.T) {
	util := util.New()
	viper.Set("jwt.refresh_token.expires_at", 1)
	viper.Set("jwt.refresh_token.secret", "secretKey")

	t.Run("When create a token, should response token and not have an error", func(t *testing.T) {
		refreshToken, err := util.JWTCreateRefreshTokenString("example@gmail.com")
		assert.Nil(t, err)
		assert.NotEmpty(t, refreshToken)
	})
}

func TestJWTCreateAccessToken(t *testing.T) {
	util := util.New()
	viper.Set("jwt.access_token.expires_at", 1)
	viper.Set("jwt.access_token.secret", "secretKey")

	t.Run("When create a token, should response token and not have an error", func(t *testing.T) {
		accessToken, err := util.JWTCreateAccessTokenString("example@gmail.com")
		assert.Nil(t, err)
		assert.NotEmpty(t, accessToken)
	})
}

func TestExtractJWTTokenFromHeaders(t *testing.T) {
	util := util.New()
	viper.Set("jwt.access_token.expires_at", time.Duration.Hours(1))
	viper.Set("jwt.access_token.secret", "secretKey")
	const accessTokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiY29tbWlzIiwiZXhwIjoxNzEzNDQ0OTQyLCJuYmYiOjE3MTM0NDQ5NDIsImlhdCI6MTcxMzQ0NDk0MiwianRpIjoiYTY4YTE2MTQtNjYzNy00MWI2LTkzOGEtYTU1YmU2ZTE2NzU3In0.a9Y0YZj1vpWbsgL7nINbW_na1pPnGN59vwKUCfPQv1c"

	t.Run("Correct Authorization header format", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {
				"Bearer " + accessTokenString,
			},
		}
		accessTokenStringFromHeader, err := util.ExtractJWTTokenStringFromHeaders(headers)
		assert.Nil(t, err)
		assert.NotEmpty(t, accessTokenStringFromHeader)
		assert.Equal(t, accessTokenString, accessTokenStringFromHeader)
	})

	t.Run("When Authorization header is empty", func(t *testing.T) {
		headers := map[string][]string{}
		accessTokenString, err := util.ExtractJWTTokenStringFromHeaders(headers)
		assert.NotNil(t, err)
		assert.Empty(t, accessTokenString)
	})

	t.Run("When Authorization header is the empty string", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {
				"",
			},
		}
		accessTokenString, err := util.ExtractJWTTokenStringFromHeaders(headers)
		assert.NotNil(t, err)
		assert.Empty(t, accessTokenString)
	})

	t.Run("When a token attached to Authorization header but it has incorrect format(token is an empty)", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {
				"Bearer",
			},
		}
		accessTokenString, err := util.ExtractJWTTokenStringFromHeaders(headers)
		assert.NotNil(t, err)
		assert.Empty(t, accessTokenString)

		headers = map[string][]string{
			"Authorization": {
				"Bearer ",
			},
		}
		accessTokenString, err = util.ExtractJWTTokenStringFromHeaders(headers)
		assert.NotNil(t, err)
		assert.Empty(t, accessTokenString)
	})

	t.Run("When a token attached to Authorization header but incorrect format(doesn’t have bearer prefix)", func(t *testing.T) {
		headers := map[string][]string{
			"Authorization": {
				accessTokenString,
			},
		}
		accessTokenString, err := util.ExtractJWTTokenStringFromHeaders(headers)
		assert.NotNil(t, err)
		assert.Empty(t, accessTokenString)
	})

}

func TestParseAccessToken(t *testing.T) {
	t.Run("", func(t *testing.T) {

	})
}

func TestParseRefreshToken(t *testing.T) {

}
