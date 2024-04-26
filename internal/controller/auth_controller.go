package controller

import (
	"booking/internal/dto"
	"booking/internal/service"

	"github.com/gofiber/fiber/v2"
)

type authController struct {
	authService service.IAuthService
}

func NewAuthController(authService service.IAuthService) *authController {
	return &authController{authService: authService}
}

func (authCtrl *authController) Register(c *fiber.Ctx) error {
	body := new(dto.RegisterDTO)
	c.BodyParser(body)
	res, statusCode := authCtrl.authService.Register(body)
	return c.Status(statusCode).JSON(res)
}

func (authCtrl *authController) Login(c *fiber.Ctx) error {
	body := new(dto.LoginDTO)
	c.BodyParser(body)
	res, statusCode := authCtrl.authService.Login(body)
	return c.Status(statusCode).JSON(res)
}

func (authCtrl *authController) RefreshToken(c *fiber.Ctx) error {
	body := new(dto.RefreshTokenDTO)
	c.BodyParser(body)
	authCtrl.authService.RefreshToken(body)
	return nil
}
