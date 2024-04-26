package service_test

import (
	"booking/internal/entity"
	"booking/internal/vo"

	"github.com/stretchr/testify/mock"
)

type authRepositoryMock struct {
	mock.Mock
}

func (repo *authRepositoryMock) CreateAccount(user *entity.User) (err error) {
	args := repo.Called(user.Username)
	return args.Error(0)
}

func (repo *authRepositoryMock) FindUserByUsername(username string) (user *entity.User, err error) {
	args := repo.Called(username)
	return args.Get(0).(*entity.User), args.Error(1)
}

type utilMock struct {
	mock.Mock
}

func (u *utilMock) JWTCreateRefreshTokenString(username string) (string, error) {
	args := u.Called(username)
	return args.String(0), args.Error(1)
}

func (u *utilMock) JWTCreateAccessTokenString(username string) (string, error) {
	args := u.Called(username)
	return args.String(0), args.Error(1)
}

func (u *utilMock) HashPassword(password string, round int) (string, error) {
	args := u.Called(password, round)
	return args.String(0), args.Error(1)
}

func (u *utilMock) ComparePassword(hashPassword string, password string) error {
	args := u.Called(hashPassword, password)
	return args.Error(0)
}

func (u *utilMock) ExtractJWTTokenStringFromHeaders(headers map[string][]string) (string, error) {
	args := u.Called(headers)
	return args.String(0), args.Error(1)
}

func (u *utilMock) JWTParseAccessToken(accessTokenString string) (*vo.AccessToken, error) {
	args := u.Called(accessTokenString)
	return args.Get(0).(*vo.AccessToken), args.Error(1)
}

func (u *utilMock) JWTParseRefreshToken(refreshTokenString string) (*vo.RefreshToken, error) {
	args := u.Called(refreshTokenString)
	return args.Get(0).(*vo.RefreshToken), args.Error(1)
}

func (u *utilMock) GetAbsoluteProjectPath() string {
	return ""
}

type cacheMock struct {
	mock.Mock
}

func (cache *cacheMock) SaveRefreshTokenString(username string, refreshTokenString string) error {
	args := cache.Called(username, refreshTokenString)
	return args.Error(0)
}

func (cache *cacheMock) SaveAccessTokenString(username string, accessTokenString string) error {
	args := cache.Called(username, accessTokenString)
	return args.Error(0)
}

func (cache *cacheMock) GetRefreshTokenString(username string) (string, error) {
	args := cache.Called(username)
	return args.String(0), args.Error(1)
}

func (cache *cacheMock) GetRefreshTokenKey(username string) string {
	args := cache.Called(username)
	return args.String(0)
}

func (cache *cacheMock) GetAccessTokenString(username string) (string, error) {
	args := cache.Called(username)
	return args.String(0), args.Error(1)
}

func (cache *cacheMock) GetAccessTokenKey(username string) string {
	args := cache.Called(username)
	return args.String(0)
}
