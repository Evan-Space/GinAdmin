package bootstrapx

import (
	"GinAdmin/config"
	"GinAdmin/data"
	casbinx "GinAdmin/internal/access/casbin"
	"GinAdmin/internal/pkg/logger"
	"GinAdmin/internal/validator"
	"fmt"
	"time"
	svc "GinAdmin/internal/service"
)

/**
* 初始化配置
**/
func InitializeConfig(configPath string) error {
	return config.InitConfig(configPath)
}

func InitializeTimezone() {
	cfg := config.GetConfig()

	timezone := "Asia/Shanghai"
	if cfg.Timezone != nil && *cfg.Timezone != "" { // 如果配置信息中时区不为空，则使用配置信息中的时区
		timezone = *cfg.Timezone
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return
	}
	time.Local = location
}

/**
* 初始化日志
 */
func InitializeLogger() error {
	return logger.InitLogger()
}

func InitializeValidator() error {
	return validator.InitValidator()
}

func Bootstrap(configPath string) error {
	if err := InitializeConfig(configPath); err != nil {
		return err
	}
	InitializeTimezone()

	if err := InitializeValidator(); err != nil {
		return err
	}

	if err := InitializeLogger(); err != nil {
		return err
	}

	if err := data.Initialize(); err != nil {
		return err
	}

	if err := casbinx.InitEnforcer(); err != nil {
		return fmt.Errorf("初始化 Casbin 失败: %w", err)
	}

	if err := svc.SyncAllPolicies(); err != nil {
		return fmt.Errorf("同步权限策略失败: %w", err)
	}
	return nil

}
