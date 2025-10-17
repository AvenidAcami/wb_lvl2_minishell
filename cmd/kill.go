package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill process by PID",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Need PID")
			return
		}
		idInt, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("Wrong PID")
			return
		}
		process, err := os.FindProcess(idInt)
		if err != nil {
			fmt.Println("Process not found:", err)
			return
		}

		if err := process.Kill(); err != nil {
			fmt.Println("Kill error:", err)
			return
		}

		_, _ = process.Wait()

	},
}
