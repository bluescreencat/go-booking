package service

import (
	"booking/internal/cache"
	"booking/internal/constant"
	"booking/internal/dto"
	"booking/internal/entity"
	"booking/internal/logs"
	"booking/internal/repository"
	"booking/internal/util"
	"booking/internal/vo"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type IAuthService interface {
	Register(registerDTO *dto.RegisterDTO) (res vo.Response, statusCode int)
	Login(loginDTO *dto.LoginDTO) (res vo.Response, statusCode int)
	RefreshToken(refreshTokenDTO *dto.RefreshTokenDTO) (res vo.Response, statusCode int)
}

type authService struct {
	cache          cache.ICache
	authRepository repository.IAuthRepository
	util           util.IUtility
}

func NewAuthService(cache cache.ICache, util util.IUtility, authRepository repository.IAuthRepository) *authService {
	return &authService{cache: cache, util: util, authRepository: authRepository}
}

func (authService *authService) Register(registerDTO *dto.RegisterDTO) (res vo.Response, statusCode int) {
	_, err := authService.authRepository.FindUserByUsername(registerDTO.Username)

	if err == nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_USERNAME_ALREADY_EXISTS)
		return res, fiber.StatusConflict
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	// create account
	hashPassword, err := authService.util.HashPassword(registerDTO.Password, viper.GetInt("password.round"))
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}
	user := entity.User{}
	user.Username = registerDTO.Username
	user.Password = hashPassword
	err = authService.authRepository.CreateAccount(&user)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	refreshTokenString, err := authService.util.JWTCreateRefreshTokenString(user.Username)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	accessTokenString, err := authService.util.JWTCreateAccessTokenString(user.Username)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	err = authService.cache.SaveRefreshTokenString(user.Username, refreshTokenString)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	err = authService.cache.SaveAccessTokenString(user.Username, accessTokenString)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	res.SetData(vo.RegisterResponse{
		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	})

	return res, fiber.StatusCreated
}

func (authService *authService) Login(loginDTO *dto.LoginDTO) (res vo.Response, statusCode int) {
	res = vo.Response{}
	user, err := authService.authRepository.FindUserByUsername(loginDTO.Username)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusUnauthorized
	}
	err = authService.util.ComparePassword(user.Password, loginDTO.Password)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusUnauthorized
	}

	refreshTokenString, err := authService.util.JWTCreateRefreshTokenString(user.Username)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	accessTokenString, err := authService.util.JWTCreateAccessTokenString(user.Username)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	err = authService.cache.SaveRefreshTokenString(user.Username, refreshTokenString)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	err = authService.cache.SaveAccessTokenString(user.Username, accessTokenString)
	if err != nil {
		logs.Error(err)
		res.SetErrorMessage(constant.ERROR_MESSAGE_INVALID_USERNAME_OR_PASSWORD)
		return res, fiber.StatusInternalServerError
	}

	res.SetData(vo.LoginResponse{
		RefreshToken: refreshTokenString,
		AccessToken:  accessTokenString,
	})

	return res, fiber.StatusOK
}

func (authService *authService) RefreshToken(refreshTokenDTO *dto.RefreshTokenDTO) (res vo.Response, statusCode int) {

	// Parse refresh token
	refreshToken, err := authService.util.JWTParseRefreshToken(refreshTokenDTO.RefreshToken)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
		return res, fiber.StatusUnauthorized
	}

	// Get refresh token from cache
	refreshTokenFromCache, err := authService.cache.GetRefreshTokenString(refreshToken.Username)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
		return res, fiber.StatusUnauthorized
	}

	// Compare the token from request with token from the cache
	if refreshTokenDTO.RefreshToken != refreshTokenFromCache {
		res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
		return res, fiber.StatusUnauthorized
	}

	// Create new access token
	accessTokenString, err := authService.util.JWTCreateAccessTokenString(refreshToken.Username)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_INTERNAL_SERVER_ERROR)
		return res, fiber.StatusInternalServerError
	}

	// Save access token
	err = authService.cache.SaveAccessTokenString(refreshToken.Username, accessTokenString)
	if err != nil {
		res.SetErrorMessage(constant.ERROR_MESSAGE_UNAUTHORIZED)
		return res, fiber.StatusInternalServerError
	}

	// Response access token
	refreshTokenResponse := vo.RefreshTokenResponse{AccessToken: accessTokenString}
	res.SetData(refreshTokenResponse)

	return res, fiber.StatusOK
}
