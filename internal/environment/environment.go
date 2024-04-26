package environment

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

func InitializeEnvironmentVariable() error {
	environment := flag.String("env", "default", "")
	flag.Parse()
	if !(*environment == "test" || *environment == "dev" || *environment == "prod") {
		return fmt.Errorf(`not yet set the environment please run with environment flag. example "go run main.go -env=dev"`)
	}
	absoluteProjectPath, _ := os.Getwd()
	absoluteProjectPath = strings.ReplaceAll(absoluteProjectPath, "\\", "/")
	extension := "yml"
	fileName := fmt.Sprintf("env_%v.%v", *environment, extension)
	viper.SetConfigFile(fileName)
	viper.AddConfigPath(absoluteProjectPath)
	fmt.Println("fileName:", fileName, "extension:", extension, "project path:", absoluteProjectPath)
	// check from environment. if it's empty will load from config file
	viper.AutomaticEnv()
	// replace . with _ for support set environment from commandline. e.g. db.port can set in commandline with parameter DB_HOST=543
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		println("err:", err.Error())
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found")
		}
	}
	return nil
}
