package vo

import "github.com/golang-jwt/jwt/v5"

type AccessToken struct {
	RefreshToken
	jwt.RegisteredClaims
}
