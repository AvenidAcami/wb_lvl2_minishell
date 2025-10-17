package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Show processes",
	Run: func(cmd *cobra.Command, args []string) {
		ps := exec.Command("ps", "aux")
		ps.Stdout = cmd.OutOrStdout()
		ps.Stderr = cmd.ErrOrStderr()
		if err := ps.Run(); err != nil {
			fmt.Println("error:", err)
		}
	},
}
