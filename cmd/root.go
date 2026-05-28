package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "Ginadmin",
	Short: "GinAdmin is a web framework for building RESTful APIs",
	Long: "GinAdmin is a web framework for building RESTful APIs",
	// Run: func(cmd *cobra.Context) {
	// 	fmt.Println("Hello, World!")
	// },
}


func Execute() {

	// err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }

}