package util

import (
	"booking/internal/constant"
	"booking/internal/vo"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func (u *utility) JWTCreateRefreshTokenString(username string) (refreshTokenString string, err error) {
	claims := vo.RefreshToken{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.refresh_token.expires_at") * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "booking",
			ID:        uuid.NewString(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshTokenString, err = token.SignedString([]byte(viper.GetString("jwt.refresh_token.secret")))
	return refreshTokenString, err
}

func (u *utility) JWTCreateAccessTokenString(username string) (accessTokenString string, err error) {
	claims := vo.AccessToken{
		RefreshToken: vo.RefreshToken{
			Username: username,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(viper.GetDuration("jwt.access_token.expires_at") * time.Hour)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				Issuer:    "booking",
				ID:        uuid.NewString(),
			},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err = token.SignedString([]byte(viper.GetString("jwt.access_token.secret")))
	return accessTokenString, err
}

func (u *utility) ExtractJWTTokenStringFromHeaders(headers map[string][]string) (accessTokenString string, err error) {

	const (
		HEADER_AUTHORIZATION_KEY        = "Authorization"
		SIZE_AUTHORIZATION_HEADER       = 1
		HEADER_JWT_TOKEN_INDEX          = 0
		SIZE_SPLIT_AUTHORIZATION_HEADER = 2
		ACCESS_TOKEN_INDEX              = 1
	)

	if len(headers[HEADER_AUTHORIZATION_KEY]) != SIZE_AUTHORIZATION_HEADER {
		return "", fmt.Errorf(constant.ERROR_MESSAGE_UNAUTHORIZED)
	}

	authorization := headers[HEADER_AUTHORIZATION_KEY][HEADER_JWT_TOKEN_INDEX]

	splitAuthorization := strings.Split(authorization, " ")
	if len(splitAuthorization) != SIZE_SPLIT_AUTHORIZATION_HEADER {
		return "", fmt.Errorf(constant.ERROR_MESSAGE_UNAUTHORIZED)
	}

	accessTokenString = splitAuthorization[ACCESS_TOKEN_INDEX]
	if strings.TrimSpace(accessTokenString) == "" {
		return "", fmt.Errorf(constant.ERROR_MESSAGE_UNAUTHORIZED)
	}
	return accessTokenString, err
}

func (u *utility) JWTParseAccessToken(accessTokenString string) (*vo.AccessToken, error) {

	secret := viper.GetString("jwt.access_token.secret")
	accessToken := vo.AccessToken{}
	token, err := jwt.ParseWithClaims(accessTokenString, &accessToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*vo.AccessToken); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims type, cannot proceed")
	}
}

func (u *utility) JWTParseRefreshToken(refreshTokenString string) (*vo.RefreshToken, error) {
	token, err := jwt.ParseWithClaims(refreshTokenString, &vo.RefreshToken{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("jwt.refresh_token.secret")), nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else if claims, ok := token.Claims.(*vo.RefreshToken); ok {
		return claims, nil
	} else {
		log.Fatal("Unknown claims type, cannot proceed")
		return nil, err
	}
}
