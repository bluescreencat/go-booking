package service_test

import (
	"booking/internal/cache"
	"booking/internal/dto"
	"booking/internal/entity"
	"booking/internal/service"
	"booking/internal/util"
	"booking/internal/vo"
	"fmt"
	"testing"
	"time"

	"github.com/XANi/loremipsum"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/storage/memory/v2"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRegister(t *testing.T) {

	authRepositoryMock := &authRepositoryMock{}
	redisMock := memory.New()
	cache := cache.New(redisMock)
	util := util.New()
	utilMock := &utilMock{}
	cacheMock := &cacheMock{}
	const ExistingEmail = "existing@gmail.com"
	const NewEmail = "newemail@gmail.com"
	mockUUID, _ := uuid.NewUUID()
	const (
		AccessTokenString  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiY29tbWlzIiwiZXhwIjoxNzEzNDQ0OTQyLCJuYmYiOjE3MTM0NDQ5NDIsImlhdCI6MTcxMzQ0NDk0MiwianRpIjoiYTY4YTE2MTQtNjYzNy00MWI2LTkzOGEtYTU1YmU2ZTE2NzU3In0.a9Y0YZj1vpWbsgL7nINbW_na1pPnGN59vwKUCfPQv1c"
		RefreshTokenString = AccessTokenString
		HashPassword       = "$2a$10$qIjfaobqHNAjvm.S0LQbYOu4W6VOwjGC2ig76PIWVi5UoM.rPZK8a"
	)

	t.Run("When get user by email and the database have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrInvalidDB)
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When email is already have, should response http status code 409 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  "dummy_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		body := dto.RegisterDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusConflict, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When have an error during a hashing password, should response error status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		lorem := loremipsum.New()
		tooLongPassword := lorem.Sentences(5)
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: tooLongPassword,
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create user and the database have an error or database timeout, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(gorm.ErrUnsupportedDriver)
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", 10)
		authService := service.NewAuthService(cache, util, authRepositoryMock)
		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create the refresh token and the function have an error, should response status code 500 with error message", func(t *testing.T) {
		const (
			Password = "dummy_password"
			Round    = 10
		)
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(nil)
		utilMock.On("HashPassword", Password, Round).Return(HashPassword, nil)
		utilMock.On("JWTCreateRefreshTokenString", NewEmail).Return("", fmt.Errorf("dummy_error"))

		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", Round)
		authService := service.NewAuthService(cache, utilMock, authRepositoryMock)
		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create the access token and the function have an error, should response status code 500 with error message", func(t *testing.T) {
		const (
			Password = "dummy_password"
			Round    = 10
		)
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(nil)
		utilMock.On("HashPassword", Password, Round).Return(HashPassword, nil)
		utilMock.On("JWTCreateRefreshTokenString", NewEmail).Return(RefreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", NewEmail).Return("", fmt.Errorf("dummy_error"))
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", Round)
		authService := service.NewAuthService(cache, utilMock, authRepositoryMock)
		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When save the refresh token and the function have an error", func(t *testing.T) {
		const (
			Password = "dummy_password"
			Round    = 10
		)
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(nil)
		utilMock.On("HashPassword", Password, Round).Return(HashPassword, nil)
		utilMock.On("JWTCreateAccessTokenString", NewEmail).Return(AccessTokenString, nil)
		utilMock.On("JWTCreateRefreshTokenString", NewEmail).Return(RefreshTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", NewEmail, RefreshTokenString).Return(fmt.Errorf("dummy_error"))
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", Round)
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		_, statusCode := authService.Register(&body)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
	})

	t.Run("When save the access token and the function have an error, should response status code 500 with error message", func(t *testing.T) {
		const (
			Password = "dummy_password"
			Round    = 10
		)
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(nil)
		utilMock.On("HashPassword", Password, Round).Return(HashPassword, nil)
		utilMock.On("JWTCreateAccessTokenString", NewEmail).Return(AccessTokenString, nil)
		utilMock.On("JWTCreateRefreshTokenString", NewEmail).Return(RefreshTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", NewEmail, RefreshTokenString).Return(nil)
		cacheMock.On("SaveAccessTokenString", NewEmail, RefreshTokenString).Return(fmt.Errorf("dummy_error"))
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", Round)
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When register is success, should return token and status 201 with token", func(t *testing.T) {
		const (
			Password = "dummy_password"
			Round    = 10
		)
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		authRepositoryMock.On("CreateAccount", NewEmail).Return(nil)
		utilMock.On("HashPassword", Password, Round).Return(HashPassword, nil)
		utilMock.On("JWTCreateAccessTokenString", NewEmail).Return(AccessTokenString, nil)
		utilMock.On("JWTCreateRefreshTokenString", NewEmail).Return(RefreshTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", NewEmail, RefreshTokenString).Return(nil)
		cacheMock.On("SaveAccessTokenString", NewEmail, RefreshTokenString).Return(nil)
		body := dto.RegisterDTO{
			Username: NewEmail,
			Password: "dummy_password",
			Name:     "dummy_name",
			Surname:  "dummy_surname",
		}
		viper.Set("password.round", Round)
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		res, statusCode := authService.Register(&body)

		isPass := assert.Equal(t, fiber.StatusCreated, statusCode)
		if isPass {
			registerResponse := res.Data.(vo.RegisterResponse)
			assert.Equal(t, AccessTokenString, registerResponse.AccessToken)
			assert.Equal(t, RefreshTokenString, registerResponse.RefreshToken)
		}
	})

}

func TestLogin(t *testing.T) {
	authRepositoryMock := &authRepositoryMock{}
	redisMock := memory.New()
	cache := cache.New(redisMock)
	util := util.New()
	utilMock := &utilMock{}
	cacheMock := &cacheMock{}
	const ExistingEmail = "existing@gmail.com"
	const NewEmail = "newemail@gmail.com"
	mockUUID, _ := uuid.NewUUID()

	const (
		AccessTokenString  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiY29tbWlzIiwiZXhwIjoxNzEzNDQ0OTQyLCJuYmYiOjE3MTM0NDQ5NDIsImlhdCI6MTcxMzQ0NDk0MiwianRpIjoiYTY4YTE2MTQtNjYzNy00MWI2LTkzOGEtYTU1YmU2ZTE2NzU3In0.a9Y0YZj1vpWbsgL7nINbW_na1pPnGN59vwKUCfPQv1c"
		RefreshTokenString = AccessTokenString
		HashPassword       = "$2a$10$qIjfaobqHNAjvm.S0LQbYOu4W6VOwjGC2ig76PIWVi5UoM.rPZK8a"
	)

	t.Run("When get user by email and the database have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{}, gorm.ErrInvalidDB)
		body := dto.LoginDTO{
			Username: NewEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When not found user in the database, should response http status code 401 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{}, gorm.ErrRecordNotFound)
		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When have an error during compare password, should response error status code 401 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", NewEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  "dummy_password",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		body := dto.LoginDTO{
			Username: NewEmail,
			Password: "dummy_hash_password",
		}
		authService := service.NewAuthService(cache, util, authRepositoryMock)

		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create refresh token and have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  HashPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		authRepositoryMock.On("CreateAccount", ExistingEmail).Return(nil)
		utilMock.On("ComparePassword", HashPassword, "dummy_password").Return(nil)
		utilMock.On("JWTCreateRefreshTokenString", ExistingEmail).Return("", fmt.Errorf("dummy_error"))

		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cache, utilMock, authRepositoryMock)
		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create access token and have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  HashPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		authRepositoryMock.On("CreateAccount", ExistingEmail).Return(nil)
		utilMock.On("ComparePassword", HashPassword, "dummy_password").Return(nil)
		utilMock.On("JWTCreateRefreshTokenString", ExistingEmail).Return(RefreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", ExistingEmail).Return("", fmt.Errorf("dummy_error"))

		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cache, utilMock, authRepositoryMock)
		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When save refresh token and have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  HashPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		authRepositoryMock.On("CreateAccount", ExistingEmail).Return(nil)
		utilMock.On("ComparePassword", HashPassword, "dummy_password").Return(nil)
		utilMock.On("JWTCreateRefreshTokenString", ExistingEmail).Return(RefreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", ExistingEmail).Return(AccessTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", ExistingEmail, RefreshTokenString).Return(fmt.Errorf("dummy_error"))

		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When save access token and have an error, should response status code 500 with error message", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  HashPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		authRepositoryMock.On("CreateAccount", ExistingEmail).Return(nil)
		utilMock.On("ComparePassword", HashPassword, "dummy_password").Return(nil)
		utilMock.On("JWTCreateRefreshTokenString", ExistingEmail).Return(RefreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", ExistingEmail).Return(AccessTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", ExistingEmail, RefreshTokenString).Return(nil)
		cacheMock.On("SaveAccessTokenString", ExistingEmail, AccessTokenString).Return(fmt.Errorf("dummy_error"))

		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When success login, should response status 200 with token", func(t *testing.T) {
		authRepositoryMock.On("FindUserByUsername", ExistingEmail).Return(&entity.User{
			ID:        mockUUID,
			Username:  ExistingEmail,
			Password:  HashPassword,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}, nil)
		authRepositoryMock.On("CreateAccount", ExistingEmail).Return(nil)
		utilMock.On("ComparePassword", HashPassword, "dummy_password").Return(nil)
		utilMock.On("JWTCreateRefreshTokenString", ExistingEmail).Return(RefreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", ExistingEmail).Return(AccessTokenString, nil)
		cacheMock.On("SaveRefreshTokenString", ExistingEmail, RefreshTokenString).Return(nil)
		cacheMock.On("SaveAccessTokenString", ExistingEmail, AccessTokenString).Return(nil)

		body := dto.LoginDTO{
			Username: ExistingEmail,
			Password: "dummy_password",
		}
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		res, statusCode := authService.Login(&body)

		isPass := assert.Equal(t, fiber.StatusOK, statusCode)
		if isPass {
			loginResponse := res.Data.(vo.LoginResponse)
			assert.Equal(t, RefreshTokenString, loginResponse.RefreshToken)
			assert.Equal(t, AccessTokenString, loginResponse.AccessToken)
		}
	})
}

func TestRefreshToken(t *testing.T) {
	util := util.New()
	authRepositoryMock := &authRepositoryMock{}
	utilMock := &utilMock{}
	cacheMock := &cacheMock{}
	const (
		Username        = "example@gmail.com"
		HashPassword    = "$2a$10$qIjfaobqHNAjvm.S0LQbYOu4W6VOwjGC2ig76PIWVi5UoM.rPZK8a"
		OldRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImV4YW1wbGVAZ21haWwuY29tIiwiaXNzIjoiY29tbWlzIiwiZXhwIjoxNzEzNDQ0OTQyLCJuYmYiOjE3MTM0NDQ5NDIsImlhdCI6MTcxMzQ0NDk0MiwianRpIjoiYTY4YTE2MTQtNjYzNy00MWI2LTkzOGEtYTU1YmU2ZTE2NzU3In0.a9Y0YZj1vpWbsgL7nINbW_na1pPnGN59vwKUCfPQv1c"
	)
	viper.Set("jwt.refresh_token.expires_at", 1*time.Hour)
	refreshTokenString, _ := util.JWTCreateRefreshTokenString(Username)
	refreshToken, _ := util.JWTParseRefreshToken(refreshTokenString)
	accessTokenString, _ := util.JWTCreateAccessTokenString(Username)

	t.Run("When parse the refresh token and the function have an error, should return status code 401 with error message", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(&vo.RefreshToken{}, fmt.Errorf("dummy_error"))
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When get refresh token from cache and the cache have an error, should return status code 401 with error message", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(refreshToken, nil)
		cacheMock.On("GetRefreshTokenString", Username).Return("", fmt.Errorf("dummy_error"))
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When compare the token from request with token from cache and these token is not the same, should return status code 401 with error message", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(refreshToken, nil)
		cacheMock.On("GetRefreshTokenString", Username).Return(OldRefreshToken, nil)
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusUnauthorized, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When create a new access token and the function have an error, should return status code 500 with error message", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(refreshToken, nil)
		cacheMock.On("GetRefreshTokenString", Username).Return(refreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", Username).Return("", fmt.Errorf("dummy_error"))
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When save access token to cache the cache have an error, should return status code 500 with error message", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(refreshToken, nil)
		cacheMock.On("GetRefreshTokenString", Username).Return(refreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", Username).Return(accessTokenString, nil)
		cacheMock.On("SaveAccessTokenString", Username, accessTokenString).Return(fmt.Errorf("dummy_error"))
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		if isPass {
			assert.NotEmpty(t, res.ErrorMessage)
		}
	})

	t.Run("When success for refresh the token, should return status 200 with new access token", func(t *testing.T) {
		utilMock.On("JWTParseRefreshToken", refreshTokenString).Return(refreshToken, nil)
		cacheMock.On("GetRefreshTokenString", Username).Return(refreshTokenString, nil)
		utilMock.On("JWTCreateAccessTokenString", Username).Return(accessTokenString, nil)
		cacheMock.On("SaveAccessTokenString", Username, accessTokenString).Return(nil)
		authService := service.NewAuthService(cacheMock, utilMock, authRepositoryMock)
		body := dto.RefreshTokenDTO{
			RefreshToken: refreshTokenString,
		}

		res, statusCode := authService.RefreshToken(&body)
		isPass := assert.Equal(t, fiber.StatusOK, statusCode)
		if isPass {
			refreshTokenResponse := res.Data.(vo.RefreshTokenResponse)
			assert.Equal(t, accessTokenString, refreshTokenResponse.AccessToken)
		}
	})
}
