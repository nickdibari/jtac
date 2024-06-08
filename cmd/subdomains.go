package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/nickdibari/jtac/src/subdomains"
	"github.com/spf13/cobra"
)

var subdomainCmd = &cobra.Command{
	Use:   "subdomains [host]",
	Short: "Lookup possible subdomains for a given host",
	Long:  `Lookup possible subdomains for a given host.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return fmt.Errorf("missing host value.")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()

		hostInput := args[0]

		fmt.Printf("performing subdomain lookup for %s...\n", hostInput)

		subdomainInfoCmd, err := subdomains.NewSubdomainInfoCmd(ctx, cmd)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		subdomains, err := subdomainInfoCmd.Run(ctx, hostInput)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		fmt.Printf("found subdomains for %s:\n", hostInput)

		for _, subdomain := range subdomains {
			fmt.Println(subdomain)
		}
	},
}

func init() {
	subdomainCmd.Flags().StringP("wordlist", "w", "", "custom wordlist to use for brute-forcing subdomains")
	rootCmd.AddCommand(subdomainCmd)
}
