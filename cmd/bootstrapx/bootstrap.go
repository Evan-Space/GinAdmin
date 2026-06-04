package bootstrapx

import (
	"fmt"
	"log"
	"time"

	"go.uber.org/zap"
)

/**
* 初始化配置
**/
func InitializeConfig(configPath string) error {
	return config.InitConfig(configPath)
}


func InitializeTimezone() {
	cfg := config.GetConfig()
	if cfg.Timezone == nil {
		return
	}

	location, err := time.LoadLocation(*cfg.Timezone)
	if err != nil {
		// 如果有错误日志记录，则记录错误日志
		if log.Logger != nil {
			log.Logger.Error(fmt.Sprintf(errorLoadingLocation, err), zap.Error(err))
		}
		fmt.Println(errorLoadingLocation+"\n", err)
		return
	}
	time.Local = location
}