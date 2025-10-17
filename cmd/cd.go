package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Change current directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("cd: missing argument")
			return
		}
		if err := os.Chdir(args[0]); err != nil {
			fmt.Println("cd error:", err)
		}
	},
}
