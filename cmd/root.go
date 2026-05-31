package cmd

import (
	"fmt"
	"os"

	"GinAdmin/cmd/service"

	"github.com/spf13/cobra"
)

/**
* 注册一个命令入口
 */
var rootCmd = &cobra.Command{
	Use: "gin-admin",
	Short: "GinAdmin is a web framework for building RESTful APIs",
	Long: "GinAdmin is a web framework for building RESTful APIss",
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


func init() {
	registerCommands()
}

func registerCommands() {
	rootCmd.AddCommand(service.ServeCmd)
}