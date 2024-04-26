package vo

import "github.com/golang-jwt/jwt/v5"

type RefreshToken struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
