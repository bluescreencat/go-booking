package config

import (
	"github.com/goccy/go-json"

	"github.com/gofiber/fiber/v2"
)

var AppConfig = fiber.Config{
	Prefork:       true,
	CaseSensitive: true,
	JSONEncoder:   json.Marshal,
	JSONDecoder:   json.Unmarshal,
}
