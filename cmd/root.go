package cmd

import (
	"fmt"
	"os"
	"strings"

	"GinAdmin/cmd/service"

	"github.com/spf13/cobra"
)


var (
	configPath string
)
/**
* 注册一个命令入口
 */
var rootCmd = &cobra.Command{
	Use: "gin-admin",
	Short: "GinAdmin is a web framework for building RESTful APIs",
	Long: "GinAdmin is a web framework for building RESTful APIss",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error { // 统一的初始化入口，避免每个子命令重复写启动代码。
		// Cobra 命令的完整生命周期是：PersistentPreRunE → PreRunE → RunE → PostRunE → PersistentPostRunE   
		if shouldSkipBootstrap(cmd) {
			return nil
		}

		// if err := bootstrapx.InitializeConfig(configPath); err != nil {
		// 	return err
		// }
		// bootstrapx.InitializeTimezone() // 日志服务还没写，后续放开注释。
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		err := service.RunServer()
		if err != nil {
			fmt.Println("Failed to start server:", err)
		}
	},
}


func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

}

/**
* 初始化
**/
func init() {
	registerCommands()
	// registerFlags()
}


/**
* 注册全局命令
**/
func registerCommands() {
	rootCmd.AddCommand(service.ServeCmd)
}

/**
* 注册全局命令行参数
**/
func registerFlags() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "", "The absolute path of the configuration file")
}


/**
* 判断是否跳过启动
* 改函数如果返回 true ,则不用启动后续的日志等服务
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