package runtime

import (
	"sync"
)




var registerOnce sync.Once

// func RegisterConfigReloadHandlers() {
// 	registerOnce.Do(func() {
// 		config.RegisterConfigReloadHandler(config.ConfigReloadHandler{
// 			Name: "data",
// 			Priority: 20,
// 			Handle: reloadData,
// 		})
// 	})
// }


// func reloadData(oldConfig, newConfig *config.Conf, diff config.ConfigDiff) error {
// 	if diff.MysqlChanged {
// 		if err := data.ReloadMysql(newConfig); err != nil {
// 			return fmt.Errorf("mysql reload failed: %w", err)
// 		}
// 		log.Logger.Info("MySQL runtime reloaded")
// 	}
// 	if diff.RedisChanged {
// 		if err := data.ReloadRedis(newConfig); err != nil {
// 			return fmt.Errorf("redis reload failed: %w", err)
// 		}
// 		log.Logger.Info("Redis runtime reloaded")
// 	}
// 	return nil
// }