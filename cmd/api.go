package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var address string

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Biu Framwork web api server",
	Long:  `Biu Framwork web api server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("api called")
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
	apiCmd.Flags().StringVarP(&address, "listen", "l", ":8080", "web server listen addr")
}
