package cmd

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)





var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		r := gin.Default()

		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		r.Run(":8080")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}