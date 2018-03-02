package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// the config structure with all configuration of the application
var (
	Port         int
	Verbose      bool
	ZipkinServer string
)

// Init intialize the config variables
func Init() {
	viper.SetDefault("PORT", 5000)

	if os.Getenv("ENVIRONMENT") == "DEV" {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			log.Fatal("error when reading the config file", err)
		}
	} else {
		viper.AutomaticEnv()
	}

	Port = viper.GetInt("PORT")
	Verbose = viper.GetBool("VERBOSE")
	ZipkinServer = viper.GetString("ZIPKIN_SERVER")
}
