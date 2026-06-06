package service

import (
	"GinAdmin/internal/routers"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

const (
	defaultHost            = "0.0.0.0"
	defaultPort            = 8080             // 默认端口
	gracefulShutdownTimout = 10 * time.Second // 关闭超时时间
)

var (
	ServeCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Starting server...")
			err := RunServer()
			if err != nil {
				fmt.Println("Failed to start server:", err)
			}

		},
	}
	host string
	port int
)

func init() {
	registerFlags()
}


/**
* 注册命令行参数
**/
func registerFlags() {
	ServeCmd.Flags().StringVarP(&host, "host", "H", defaultHost, "监听服务器地址")
	ServeCmd.Flags().IntVarP(&port, "port", "P", defaultPort, "监听服务器端口")
}


/**
* 启动服务器
**/
func RunServer() error {
	engine, err := routers.SetRouters()
	if err != nil {
		return fmt.Errorf("Failed to set routers: %v", err)
	}

	address := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{
		Addr:    address,
		Handler: engine,
	}
	// errChan := make(chan error, 1)
	// go func() {
	// 	// log.Logger.Info("API service starting", zap.String("address", address))
	// 	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		errChan <- err
	// 	}
	// 	close(errChan)
	// }()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	// return waitForShutdown(server, errChan)
	return nil
}