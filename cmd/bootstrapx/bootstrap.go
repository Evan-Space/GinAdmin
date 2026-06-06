package bootstrapx

import (
	"GinAdmin/config"
	"time"
)

/**
* 初始化配置
**/
func InitializeConfig(configPath string) error {
	return config.InitConfig(configPath)
}


func InitializeTimezone() error {
	cfg := config.GetConfig()

	timezone := "Asia/Shanghai"
	if cfg.Timezone != nil && *cfg.Timezone != "" { // 如果配置信息中时区不为空，则使用配置信息中的时区
		timezone = *cfg.Timezone
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return err
	}
	
	// if err != nil {
	// 	// 如果有错误日志记录，则记录错误日志
	// 	if log.Logger != nil {
	// 		log.Logger.Error(fmt.Sprintf(errorLoadingLocation, err), zap.Error(err))
	// 	}
	// 	fmt.Println(errorLoadingLocation+"\n", err)
	// 	return
	// }
	time.Local = location
	return nil
}