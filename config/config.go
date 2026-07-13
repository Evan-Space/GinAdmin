package config

import (
	"GinAdmin/config/autoload"
	"os"

	"github.com/goccy/go-yaml"
)

/**
* 整个项目配置信息的类型
**/
type Conf struct {
	JWT      JWTConfig             `yaml:"jwt"`
	MySQl    MySQLConfig           `yaml:"mysql"`
	Redis    RedisConfig           `yaml:"redis"`
	Timezone *string               `yaml:"timezone"`
	Logger   autoload.LoggerConfig `yaml:"logger"`
	BasePath string                `yaml:"base_path"`
}

type JWTConfig struct {
	SecretKey  string `yaml:"secret_key"`
	TTL        string `yaml:"ttl"`
	RefreshTTL string `yaml:"refresh_ttl` // refresh 模式下新 token 有效期（可选，比 ttl 长）
}

type MySQLConfig struct {
	Enable   bool   `yaml:"enable"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
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

	basePath, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg.BasePath = basePath

	return nil
}

func GetConfig() *Conf {
	return &cfg
}

/**
* 配置重载处理器
**/
// type ConfigReloadHandler struct {
// 	Name     string
// 	Priority int
// 	Handle   func(oldConfig, newConfig *Conf, diff ConfigDiff) error
// }
// type ConfigDiff struct {
// 	LoggerChanged         bool
// 	MysqlChanged          bool
// 	RedisChanged          bool
// 	JWTChanged            bool
// 	JWTSecretChanged      bool
// 	BaseURLChanged        bool
// 	CORSChanged           bool
// 	TrustedProxiesChanged bool
// 	LightAppChanged       bool
// 	RestartRequiredFields []string
// 	ChangedFields         []string
// }

// RegisterConfigReloadHandler 注册配置热更新回调。
// func RegisterConfigReloadHandler(handler ConfigReloadHandler) {
// if handler.Name == "" {
// 	return
// }

// reloadHandlersMu.Lock()
// defer reloadHandlersMu.Unlock()

// for i := range reloadHandlers {
// 	if reloadHandlers[i].Name == handler.Name {
// 		reloadHandlers[i] = handler
// 		sortConfigReloadHandlersLocked()
// 		return
// 	}
// }

// reloadHandlers = append(reloadHandlers, handler)
// sortConfigReloadHandlersLocked()
// }
