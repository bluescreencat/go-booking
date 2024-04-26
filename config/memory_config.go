package config

import (
	"time"

	"github.com/gofiber/storage/memory/v2"
)

var MemoryConfig = memory.Config{
	GCInterval: 60 * time.Second,
}
