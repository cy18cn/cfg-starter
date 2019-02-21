package zlog

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)


var zLog *zap.Logger

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	// viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
	// viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
	viper.AddConfigPath(".")    // optionally look for config in the working directory
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.SetDefault("logLevel", "INFO")
	viper.SetDefault("logFile", "app.log")
	viper.SetDefault("serviceName", "myApp")

	zLog = newLJZapLogger()
}

// type Config struct {
// 	DBAddr string `yaml:"dbAddr"`
// 	LoggerConfig
// }

// type LoggerConfig struct {
// 	LoggerLevel string `yaml:"level"`
// 	Filename    string `yaml:"filename"`
// 	Development bool   `yaml:"development"`
// 	ServiceName string `yaml:"serviceName"`
// }

// type DSConfig struct {
// 	DBAddr string `yaml:"dbAddr"`
// }

// var AppConfig *Config

// func LoadAll(file string) ([]*DSConfig, error) {
// 	f, err := os.Open(file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer f.Close()

// 	dec := yaml.NewDecoder(f)

// 	var configs []*DSConfig
// 	var config DSConfig
// 	for dec.Decode(&config) == nil {
// 		configs = append(configs, &config)
// 	}

// 	return configs, nil
// }

// func LoadConfig(file string) (*DSConfig, error) {
// 	configs, err := LoadAll(file)
// 	if err != nil {

// 	}
// }
