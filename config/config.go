package config

import (
	"os"

	"github.com/goccy/go-yaml"
)

/**
* 整个项目配置信息的类型
**/
type Conf struct {
	JWT JWTConfig `yaml:"jwt"`
	MySQl MySQLConfig `yaml:"mysql"`
	Redis RedisConfig `yaml:"redis"`
	Timezone *string `yaml:"timezone"`
}


type JWTConfig struct {
	SecretKey string `yaml:"secret_key"`
}


type MySQLConfig struct {
	Enable bool `yaml:"enable"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}


type RedisConfig struct {
	Enable bool `yaml:"enable"`
}

var cfg Conf 

/**
* 获取初始化配置
**/
func InitConfig(path string) error {
	if path == "" {
		path = "config/config.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}


	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return nil
}

func GetConfig() *Conf {
	return &cfg
}