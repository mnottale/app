package main

import (
	"github.com/docker/cli/cli/command"
	"github.com/spf13/cobra"
)

func completionCmd(dockerCli command.Cli, rootCmd *cobra.Command) *cobra.Command {
	return &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(docker-app completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(docker-app completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(dockerCli.Out())
		},
	}
}
