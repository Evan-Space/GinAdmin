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

func registerFlags() {
	ServeCmd.Flags().StringVarP(&host, "host", "H", defaultHost, "监听服务器地址")
	ServeCmd.Flags().IntVarP(&port, "port", "P", defaultPort, "监听服务器端口")
}


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
	errChan := make(chan error, 1)
	go func() {
		// log.Logger.Info("API service starting", zap.String("address", address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()

	return server.ListenAndServe()

	// return waitForShutdown(server, errChan)
}



// func waitForShutdown(server *http.Server, errChan <-chan error) error {
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
// 	defer signal.Stop(sigChan)

// 	select {
// 	case err, ok := <-errChan:
// 		if ok && err != nil {
// 			return err
// 		}
// 		return nil
// 	case sig := <-sigChan:
// 		log.Logger.Warn("Received API shutdown signal", zap.String("signal", sig.String()))
// 	}

// 	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimout)
// 	defer cancel()

// 	if err := server.Shutdown(ctx); err != nil {
// 		return fmt.Errorf("shutdown http server failed: %w", err)
// 	}
// 	if err := data.Shutdown(); err != nil {
// 		return fmt.Errorf("shutdown data resources failed: %w", err)
// 	}

// 	log.Logger.Info("API service stopped gracefully")
// 	return nil
// }

// package service

// import (
// 	"github.com/gin-gonic/gin"
// 	"github.com/spf13/cobra"
// )

// var serverCmd = &cobra.Command{
// 	Use: "server",
// 	Short: "Start the server",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		r := gin.Default()

// 		r.GET("/ping", func(c *gin.Context) {
// 			c.JSON(200, gin.H{
// 				"message": "pong",
// 			})
// 		})

// 		r.Run(":8080")
// 	},
// }

// func init() {
// 	rootCmd.AddCommand(serverCmd)
// }
