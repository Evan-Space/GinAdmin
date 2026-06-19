package data

import (
	"GinAdmin/config"
	"fmt"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)


var (
	db *gorm.DB // 全局数据库实例
	dbMu sync.Mutex // 初始化锁
	dbErr error // 初始化错误
)

// Initialize 初始化 MySQL 连接
func Initialize() error {
	cfg := config.GetConfig() // config 已经初始化完毕，这里可以直接获取 config 配置信息

	if !cfg.MySQl.Enable { // 如果配置信息中数据库没有开启，则直接返回
		return nil
	}

	dbMu.Lock()
	defer dbMu.Unlock()

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQl.Username,
		cfg.MySQl.Password,
		cfg.MySQl.Host,
		cfg.MySQl.Port,
		cfg.MySQl.Database,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		dbErr = err
		return err
	}


	db = conn
	return nil
}


func GetDB() *gorm.DB {
	return db
}

func Shutdown() error {
	if db == nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}


// SetDB 设置数据库实例（用于测试或事务场景注入）
func SetDB(conn *gorm.DB) {
	db = conn
}

