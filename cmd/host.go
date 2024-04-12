package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/nickdibari/jtac/src/host"
	"github.com/spf13/cobra"
)

var hostCmd = &cobra.Command{
	Use:   "host [host]",
	Short: "Get information about a particular host",
	Long: `Get information about a particular host.

This command will work with either an IP address or hostname.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("missing host value.")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		hostInput := args[0]

		fmt.Printf("performing host lookup for %s...\n", hostInput)

		hostInfoCmd, err := host.NewHostInfoCmd(ctx, cmd)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		hostInfo, err := hostInfoCmd.Run(ctx, hostInput)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("found host info for %s:\n", hostInput)
		fmt.Printf("%s\n", hostInfo.Format())
	},
}

func init() {
	hostCmd.Flags().StringP("resolver", "r", "", "custom resolver IP address to use for DNS lookups")

	rootCmd.AddCommand(hostCmd)
}
