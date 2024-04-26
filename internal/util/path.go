package util

import (
	"os"
	"strings"
)

func (u *utility) GetAbsoluteProjectPath() string {
	dir, _ := os.Getwd()
	projectPath, _, _ := strings.Cut(dir, "internal")
	return projectPath
}
