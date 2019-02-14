package cmd

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-bind9/src/version"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display this build's version, build time, and git hash",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(version.FormattedMessage())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
