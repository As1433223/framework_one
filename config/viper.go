package config

import "github.com/spf13/viper"

type NaCosConf struct {
	Ip         string `yaml:"Ip"`
	Port       int    `yaml:"Port"`
	ServerName string `yaml:"ServerName"`
}

func InitViper(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	return err
}
