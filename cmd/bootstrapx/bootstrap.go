package bootstrapx

import (
	"GinAdmin/config"
	"GinAdmin/internal/validator"
	"time"
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
func InitializeLogger()  error {
	return nil
}



func InitializeValidator() error {
	return validator.InitValidator()
}