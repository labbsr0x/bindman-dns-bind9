package nsupdate

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBingFlags(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	address := "bind"
	port := "8080"
	keyFile := "Ktest.com.+157+50086.key"
	zone := "test.com"

	err := command.ParseFlags([]string{
		fmt.Sprintf("--%s=%s", nameServerAddress, address),
		fmt.Sprintf("--%s=%s", nameServerPort, port),
		fmt.Sprintf("--%s=%s", nameServerKeyFile, keyFile),
		fmt.Sprintf("--%s=%s", nameServerZone, zone),
		fmt.Sprintf("--%s=%t", debug, true),
	})
	require.NoError(t, err)

	b := &Builder{}
	b.InitFromViper(v)

	assert.Equal(t, address, b.Server)
	assert.Equal(t, port, b.Port)
	assert.Equal(t, keyFile, b.KeyFile)
	assert.Equal(t, zone, b.Zone)
	assert.Equal(t, true, b.Debug)
}

func TestDefaultValues(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	b := &Builder{}
	b.InitFromViper(v)

	assert.Equal(t, defaultNameServerPort, b.Port)
	assert.Equal(t, false, b.Debug)
}
