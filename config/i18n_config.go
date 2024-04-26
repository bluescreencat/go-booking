package config

import (
	"path"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

func getAbsoluteLanguagePath() string {
	languageFolderName := "localize"
	projectPath := viper.GetString("project_directory")
	return path.Join(projectPath, languageFolderName)
}

var I18nConfig = &fiberi18n.Config{
	RootPath:        getAbsoluteLanguagePath(),
	AcceptLanguages: []language.Tag{language.Thai, language.English},
	DefaultLanguage: language.English,
}
