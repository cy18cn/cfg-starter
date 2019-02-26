package common

import (
	"fmt"
	"github.com/cy18cn/micro-svc-common/zlog"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
)

func InitConfig(reloadConfig func(in fsnotify.Event)) error {
	viper.SetConfigName("config")                 // name of config file (without extension)
	viper.AddConfigPath(os.Getenv("CONFIG_PATH")) // call multiple times to add many search paths
	viper.AddConfigPath("/etc/appname/")          // path to look for the config file in
	viper.AddConfigPath(".")                      // optionally look for config in the working directory
	err := viper.ReadInConfig()                   // Find and read the config file
	if err != nil {                               // Handle errors reading the config file
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	viper.SetDefault("logLevel", "INFO")
	viper.SetDefault("logFile", "app.log")
	viper.SetDefault("app.name", "myApp")

	err = zlog.InitProduction(viper.GetString("app.name"))
	if err != nil {
		return fmt.Errorf("Fatal error config file: %s \n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(reloadConfig)
	return nil
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
