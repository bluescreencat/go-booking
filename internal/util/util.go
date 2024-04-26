package util

import "booking/internal/vo"

type utility struct {
}

type IUtility interface {
	HashPassword(password string, round int) (string, error)
	ComparePassword(hashPassword string, password string) error
	JWTCreateRefreshTokenString(username string) (string, error)
	JWTCreateAccessTokenString(username string) (string, error)
	ExtractJWTTokenStringFromHeaders(headers map[string][]string) (string, error)
	JWTParseAccessToken(accessTokenString string) (*vo.AccessToken, error)
	JWTParseRefreshToken(refreshTokenString string) (*vo.RefreshToken, error)
	GetAbsoluteProjectPath() string
}

func New() *utility {
	return &utility{}
}
