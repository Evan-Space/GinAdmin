package config

import (
	"GinAdmin/config/autoload"
	"os"
	"path/filepath"
	"strings"
	"time"

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
	Storage  StorageConfig         `yaml:"storage"`
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

// 存储配置
type StorageConfig struct {
	RootDir     string   `yaml:"root_dir"`      // 存储的跟路径，相对路径是，基于 BasePath
	ChunkSize   int64    `yaml:"chunk_size"`    // 分片大小，必须与前端 CHUNK_SIZE 一致
	MaxFileSize int64    `yaml:"max_file_size"` // 文件大小上限
	AllowExt    []string `yaml:"allow_ext"`     // 扩展名白名单，为空表示不限制
	TmpTTL      string   `yaml:"tmp_ttl"`       // 未完成任务的保留时长
}

const (
	defaultChunkSize   int64 = 5 << 20        // 5MB
	defaultMaxFileSize int64 = 10 << 30       // 10GB
	defaultTmpTTL            = 24 * time.Hour // 24小时
)

/*
获取存储根路径
*/
func (c *Conf) StorageRoot() string {
	dir := c.Storage.RootDir
	if dir == "" {
		dir = "storage/uploads"
	}
	if filepath.IsAbs(dir) {
		return dir
	}

	return filepath.Join(c.BasePath, dir)
}

func (c *Conf) ChunkSize() int64 {
	if c.Storage.ChunkSize > 0 {
		return c.Storage.ChunkSize
	}
	return defaultChunkSize
}

func (c *Conf) MaxFileSize() int64 {
	if c.Storage.MaxFileSize > 0 {
		return c.Storage.MaxFileSize
	}
	return defaultMaxFileSize
}

// ExtAllowed 扩展名白名单校验
func (c *Conf) ExtAllowed(ext string) bool {
	if len(c.Storage.AllowExt) == 0 {
		return true
	}
	for _, item := range c.Storage.AllowExt {
		if strings.EqualFold(item, ext) {
			return true
		}
	}
	return false
}

// TmpTTLDuration 未完成任务的保留时长
func (c *Conf) TmpTTLDuration() time.Duration {
	if d, err := time.ParseDuration(c.Storage.TmpTTL); err == nil && d > 0 {
		return d
	}
	return defaultTmpTTL
}

// ------------------------------------------------------------
// 配置文件解析
// ------------------------------------------------------------
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
