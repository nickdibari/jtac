package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jtac",
	Short: "Tool for gathering OSINT for IP addresses and hostnames",
	Long:  `Tool for gathering OSINT for IP addresses and hostnames`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
