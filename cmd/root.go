package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "minishell",
	Short: "minishell",
	Run: func(cmd *cobra.Command, args []string) {
		runShell(cmd)
	},
}

func Execute() {
	cobra.OnInitialize(initConfig)
	rootCmd.AddCommand(pwdCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.AutomaticEnv()
}

func runShell(root *cobra.Command) {
	reader := bufio.NewScanner(os.Stdin)

	for {
		dir, _ := os.Getwd()
		fmt.Printf("%s > ", dir)

		if !reader.Scan() {
			break
		}

		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		command := args[0]
		cmdArgs := args[1:]

		c, _, err := root.Find(args)
		if err == nil && c != nil && c != root {
			fmt.Println("in comma:", command, "args: ", cmdArgs) // delete it later
			if c.Run != nil {
				c.Run(c, cmdArgs)
			}
			continue
		} else {
			fmt.Println("ext comma:", command, "args: ", cmdArgs) // delete it later
			execExternal(command, cmdArgs)
		}
	}
}

func execExternal(name string, args []string) {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		fmt.Println("error:", err)
	}
}
