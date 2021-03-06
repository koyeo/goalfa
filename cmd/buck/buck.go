package main

import (
	"github.com/koyeo/buck/cmd/buck/sdk"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "buck",
	Long: `用于辅助构建 Buck SDK 的命令行工具`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func main() {
	rootCmd.AddCommand(
		sdk.Command,
	)
	if err := rootCmd.Execute(); err != nil {
		return
	}
}
