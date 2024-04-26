package controller_test

import (
	"booking/internal/controller"
	"booking/internal/dto"
	"booking/internal/vo"
	"booking/middleware/validator"
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type authServiceMock struct {
	mock.Mock
}

func (authService *authServiceMock) Login(loginDTO *dto.LoginDTO) (res vo.Response, statusCode int) {
	return res, fiber.StatusOK
}

func (authService *authServiceMock) RefreshToken(refreshTokenDTO *dto.RefreshTokenDTO) (res vo.Response, statusCode int) {
	return res, fiber.StatusOK
}

func (authService *authServiceMock) Register(registerDTO *dto.RegisterDTO) (res vo.Response, statusCode int) {
	return res, fiber.StatusCreated
}

type TestCase[T interface{}, E interface{}] struct {
	Name     string
	Input    T
	Expected E
}

func TestRegister(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")
	auth := api.Group("/auth")
	authService := authServiceMock{}
	authCtrl := controller.NewAuthController(&authService)
	auth.Post("/register", validator.BodyValidator[dto.RegisterDTO](), authCtrl.Register)
	const URL = "/api/auth/register"

	testCases := []TestCase[dto.RegisterDTO, int]{
		{
			Name: "Valid username, password, name, surname format",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
				Name:     "jhon",
				Surname:  "smith",
			},
			Expected: fiber.StatusCreated,
		},
		{
			Name: "Invalid username format",
			Input: dto.RegisterDTO{
				Username: "not_email_format",
				Password: "aaaAAA#1",
				Name:     "jhon",
				Surname:  "smith",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only lowercase)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaaaaaa",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only uppercase)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "AAAAAAAA",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only special character)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "########",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only number)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "11111111",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (too short password)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aA#1",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (too long password)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaaaAAAAA#####12345",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid name format (nil)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
				Surname:  "smith",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid name format (empty string)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
				Name:     "",
				Surname:  "smith",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid surname format (nil)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
				Name:     "jhon",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid surname format (empty string)",
			Input: dto.RegisterDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
				Name:     "jhon",
				Surname:  "",
			},
			Expected: fiber.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input, _ := json.Marshal(testCase.Input)
			req := httptest.NewRequest(fiber.MethodPost, URL, bytes.NewReader(input))
			req.Header.Set("Content-Type", "application/json")
			res, _ := app.Test(req)
			assert.Equal(t, testCase.Expected, res.StatusCode)
		})
	}

}

func TestLogin(t *testing.T) {
	app := fiber.New()
	api := app.Group("/api")
	auth := api.Group("/auth")
	authService := authServiceMock{}
	authCtrl := controller.NewAuthController(&authService)
	auth.Post("/login", validator.BodyValidator[dto.LoginDTO](), authCtrl.Login)
	const url = "/api/auth/login"

	testCases := []TestCase[dto.LoginDTO, int]{
		{
			Name: "Valid username and password format",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "aaaAAA#1",
			},
			Expected: fiber.StatusOK,
		},
		{
			Name: "Invalid username format",
			Input: dto.LoginDTO{
				Username: "not_an_email",
				Password: "aaaAAA#1",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid username format",
			Input: dto.LoginDTO{
				Username: "not_an_email@g",
				Password: "aaaAAA#1",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only lowercase)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "aaaaaaaa",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only uppercase)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "AAAAAAAA",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only special character)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "########",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (only number)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "11111111",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (too short password)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "aA#1",
			},
			Expected: fiber.StatusBadRequest,
		},
		{
			Name: "Invalid password format (too long password)",
			Input: dto.LoginDTO{
				Username: "example@gmail.com",
				Password: "aaaaaAAAAA#####12345",
			},
			Expected: fiber.StatusBadRequest,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			input, _ := json.Marshal(testCase.Input)
			req := httptest.NewRequest(fiber.MethodPost, url, bytes.NewReader(input))
			req.Header.Set("Content-Type", "application/json")
			res, _ := app.Test(req)
			assert.Equal(t, testCase.Expected, res.StatusCode)
		})
	}

}
