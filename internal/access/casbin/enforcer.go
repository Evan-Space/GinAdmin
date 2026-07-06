package casbinx

import (
	"GinAdmin/data"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/casbin/casbin/v3"
	"github.com/casbin/casbin/v3/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

type CasbinEnforcer struct {
	*casbin.Enforcer
}

var (
	enforcer *CasbinEnforcer
	once     sync.Once
	initErr  error
)

// InitEnforcer 初始化 Casbin Enforcer（只执行一次）
// 启动时 在 data.Initialize 中调用
func InitEnforcer() error {
	once.Do(func() {
		initErr = doInit()
	})
	return initErr
}

func doInit() error {
	// 1. 先火球模型文件路径
	cwd, _ := os.Getwd()
	modelPath := filepath.Join(cwd, "rbac_model.conf")
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return fmt.Errorf("Casbin 模型文件不存在: %s", modelPath)
	}

	// 2. 加载模型
	m, err := model.NewModelFromFile(modelPath);
	if err != nil {
		return fmt.Errorf("加载 Casbin 模型失败: %w", err)
	}

	// 3. 获取数据库链接
	db := data.GetDB()
	if db == nil {
		return fmt.Errorf("数据库未初始化")
	}
	// 4. 创建 GORM 适配器（策略存储在 casbin_rule 表）
	// gormadapter.TurnOffAutoMigrate(db) // 不使用自动迁移
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return fmt.Errorf("创建 Casbin 适配器失败: %w", err)
	}

	// 5. 创建 Enforcer（自动从 casbin_rule 表加载策略）
	e, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return fmt.Errorf("创建 Enforcer 失败: %w", err)
	}
	e.EnableAutoSave(true) // 调用 AddPolicy 等方法时自动写入数据库

	enforcer = &CasbinEnforcer{Enforcer: e}
	return nil

}

// GetEnforcer 返回已初始化的 Enforcer 实例
func GetEnforcer() (*CasbinEnforcer, error) {
	if enforcer == nil {
		if err := InitEnforcer(); err != nil {
			return nil, err
		}
	}
	return enforcer, nil
}

// ReloadPolicy 重新从数据库加载策略
func ReloadPolicy() error {
	e, err := GetEnforcer()
	if err != nil {
		return err
	}
	return e.LoadPolicy()
}
