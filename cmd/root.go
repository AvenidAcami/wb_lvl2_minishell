package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
	rootCmd.AddCommand(pwdCmd, cdCmd, echoCmd, killCmd, psCmd)
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

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	var currentCmd *exec.Cmd

	go func() {
		for range sigChan {
			if currentCmd != nil && currentCmd.Process != nil {
				_ = currentCmd.Process.Kill()
				fmt.Println("\nCommand aborted by Ctrl+C")
			}
		}
	}()

	for {
		dir, _ := os.Getwd()
		fmt.Printf("%s > ", dir)

		if !reader.Scan() {
			fmt.Println("\nexit")
			break
		}

		line := strings.TrimSpace(reader.Text())
		if line == "" {
			continue
		}

		if strings.Contains(line, "|") {
			if err := runPipeline(line); err != nil {
				fmt.Println("Pipeline error:", err)
			}
			continue
		}

		args := strings.Fields(line)
		command := args[0]
		cmdArgs := args[1:]

		c, _, err := root.Find(args)
		if err == nil && c != nil && c != root {
			if c.Run != nil {
				c.Run(c, cmdArgs)
			}
			continue
		}

		cmd := exec.Command(command, cmdArgs...)
		currentCmd = cmd
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			if strings.Contains(err.Error(), "signal: interrupt") {
				continue
			}
			fmt.Println("Error:", err)
		}
		currentCmd = nil
	}
}

func runPipeline(line string) error {
	parts := strings.Split(line, "|")
	var cmds []*exec.Cmd

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		args := strings.Fields(part)
		if args[0] == "cd" {
			fmt.Println("Cant use cd in pipeline")
			return nil
		}
		cmds = append(cmds, exec.Command(args[0], args[1:]...))
	}

	if len(cmds) == 0 {
		return nil
	}

	for i := 0; i < len(cmds)-1; i++ {
		out, err := cmds[i].StdoutPipe()
		if err != nil {
			return err
		}
		cmds[i+1].Stdin = out
	}

	cmds[len(cmds)-1].Stdout = os.Stdout
	cmds[len(cmds)-1].Stderr = os.Stderr

	for _, c := range cmds {
		if err := c.Start(); err != nil {
			return err
		}
	}

	for _, c := range cmds {
		_ = c.Wait()
	}

	return nil
}
