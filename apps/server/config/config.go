package config

import (
	"GinAdmin/config/autoload"
	"os"
	"strconv"

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
	TLS      bool   `yaml:"tls"`
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
		path = os.Getenv("CONFIG_PATH")
	}
	if path == "" {
		path = "config/config.yaml"
	}

	// 配置文件缺失时不报错，完全依赖环境变量
	if data, err := os.ReadFile(path); err == nil {
		if err := yaml.Unmarshal(data, &cfg); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	applyEnv()
	basePath, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg.BasePath = basePath
	return nil

	// if path == "" {
	// 	path = "config/config.yaml"
	// }

	// data, err := os.ReadFile(path)
	// if err != nil {
	// 	return err
	// }

	// if err := yaml.Unmarshal(data, &cfg); err != nil {
	// 	return err
	// }

	// basePath, err := os.Getwd()
	// if err != nil {
	// 	return err
	// }
	// cfg.BasePath = basePath

	// return nil
}

func GetConfig() *Conf {
	return &cfg
}

// applyEnv 部署环境的密钥只放环境变量，不进代码仓库
func applyEnv() {
	setString(&cfg.JWT.SecretKey, "JWT_SECRET")
	setString(&cfg.JWT.TTL, "JWT_TTL")
	setString(&cfg.JWT.RefreshTTL, "JWT_REFRESH_TTL")
	setString(&cfg.MySQl.Host, "MYSQL_HOST")
	setString(&cfg.MySQl.Database, "MYSQL_DATABASE")
	setString(&cfg.MySQl.Username, "MYSQL_USER")
	setString(&cfg.MySQl.Password, "MYSQL_PASSWORD")
	setString(&cfg.Logger.Output, "LOG_OUTPUT")
	if port, err := strconv.Atoi(os.Getenv("MYSQL_PORT")); err == nil {
		cfg.MySQl.Port = port
	}
	if v := os.Getenv("MYSQL_TLS"); v != "" {
		cfg.MySQl.TLS = v == "true"
	}
	if os.Getenv("MYSQL_HOST") != "" {
		cfg.MySQl.Enable = true
	}
}
func setString(target *string, key string) {
	if v := os.Getenv(key); v != "" {
		*target = v
	}
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
