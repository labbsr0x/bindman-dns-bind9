package cmd

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-bind9/src/manager"
	"github.com/labbsr0x/bindman-dns-bind9/src/nsupdate"
	"github.com/labbsr0x/bindman-dns-webhook/src/hook"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const basePath = "/data"

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server and serves the HTTP REST API",
	RunE:  runE,
}

func runE(_ *cobra.Command, _ []string) error {
	nsupdateBuilder := new(nsupdate.Builder).InitFromViper(viper.GetViper())
	managerBuilder := new(manager.Builder).InitFromViper(viper.GetViper())
	nsu, err := nsupdateBuilder.New(basePath)
	if err != nil {
		return fmt.Errorf("an error occurred while setting up the DNS Manager: %v", err)
	}
	bind9Manager, err := managerBuilder.New(nsu, basePath)
	if err != nil {
		return err
	}
	hook.Initialize(bind9Manager)
	return nil
}

func init() {
	rootCmd.AddCommand(serveCmd)

	nsupdate.AddFlags(serveCmd.Flags())
	manager.AddFlags(serveCmd.Flags())

	err := viper.GetViper().BindPFlags(serveCmd.Flags())
	if err != nil {
		panic(err)
	}
}
