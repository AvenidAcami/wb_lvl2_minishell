package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var echoCmd = &cobra.Command{
	Use:   "echo",
	Short: "Print args",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(strings.Join(args, " "))
	},
}
