package app

import (
	"booking/config"
	"booking/internal/cache"
	"booking/internal/controller"
	"booking/internal/dto"
	"booking/internal/repository"
	"booking/internal/service"
	"booking/internal/util"
	"booking/middleware/validator"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

var app *fiber.App

func GetApp(cache cache.ICache, db *gorm.DB) *fiber.App {

	if app != nil {
		return app
	}

	app = fiber.New(config.AppConfig)
	app.Use(recover.New())
	app.Use(fiberi18n.New(config.I18nConfig))
	app.Use(compress.New(config.CompressConfig))
	app.Use(requestid.New())
	app.Use(cors.New())

	api := app.Group("/api")
	api.Use(logger.New(config.LoggerConfig))

	v1 := api.Group("/v1")

	auth := v1.Group("/auth")

	util := util.New()

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(cache, util, authRepository)
	authCtrl := controller.NewAuthController(authService)
	auth.Post("/register", validator.BodyValidator[dto.RegisterDTO](), authCtrl.Register)
	auth.Post("/login", validator.BodyValidator[dto.LoginDTO](), authCtrl.Login)
	auth.Post("/refresh-token", validator.BodyValidator[dto.RefreshTokenDTO](), authCtrl.RefreshToken)
	return app
}
