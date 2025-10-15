package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "Get current directory",
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(dir)
	},
}
