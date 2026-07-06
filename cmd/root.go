package cmd

import (
	"fmt"
	"os"
	"strings"

	"GinAdmin/cmd/bootstrapx"
	cmd_service "GinAdmin/cmd/service"
	"GinAdmin/data"
	svc "GinAdmin/internal/service"

	casbinx "GinAdmin/internal/access/casbin"

	"github.com/spf13/cobra"
)

var (
	configPath     string
	welcomeMessage = "Welcome to go-layout. Use -h to see more commands"
)

/**
* 注册一个命令入口
 */
var rootCmd = &cobra.Command{
	Use:   "gin-admin",
	Short: "GinAdmin is a web framework for building RESTful APIs",
	Long:  "GinAdmin is a web framework for building RESTful APIss",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { // 统一的初始化入口，避免每个子命令重复写启动代码。
		// Cobra 命令的完整生命周期是：PersistentPreRunE → PreRunE → RunE → PostRunE → PersistentPostRunE
		if shouldSkipBootstrap(cmd) {
			return nil
		}

		/**
		* 启动项目前，先初始化配置，查询配置文件，并初始化配置信息
		 */
		if err := bootstrapx.InitializeConfig(configPath); err != nil {
			return err
		}

		/**
		* 初始化时区
		 */
		bootstrapx.InitializeTimezone()
		if err := bootstrapx.InitializeValidator(); err != nil {
			return err
		}

		/**
		* 初始化数据库
		 */
		err := data.Initialize()
		if err != nil {
			return err
		}
		/**
		* 初始化 Casbin
		 */
		if err := casbinx.InitEnforcer(); err != nil {
			return fmt.Errorf("初始化 Casbin 失败: %w", err)
		}

		// 同步权限策略
		if err := svc.SyncAllPolicies(); err != nil {
			return fmt.Errorf("同步权限策略失败: %w", err)
		}

		return nil // 如果初始化数据库成功，则返回 nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n", welcomeMessage)
		// err := service.RunServer()
		// if err != nil {
		// 	fmt.Println("Failed to start server:", err)
		// }
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

/**
* 初始化
**/
func init() {
	registerFlags()
	registerCommands()
}

/**
* 注册全局命令
**/
func registerCommands() {
	rootCmd.AddCommand(cmd_service.ServeCmd)
}

/**
* 注册全局命令行参数
**/
func registerFlags() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "The absolute path of the configuration file")
}

/**
* 判断是否跳过启动
* 该函数如果返回 true ,则不用启动后续的日志等服务
**/
func shouldSkipBootstrap(cmd *cobra.Command) bool {
	if cmd == nil {
		return false
	}
	if cmd.Name() == "help" || cmd.Name() == "version" {
		return true
	}
	commandPath := cmd.CommandPath()

	switch commandPath {
	case "gin-admin", "gin-admin version", "gin-admin help":
		return true
	default:
		return strings.HasPrefix(commandPath, "gin-admin completion") || strings.HasPrefix(commandPath, "gin-admin __complete")
	}
}
