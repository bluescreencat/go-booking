package config

import "github.com/gofiber/fiber/v2/middleware/compress"

var CompressConfig = compress.Config{
	Level: compress.LevelBestSpeed,
}
