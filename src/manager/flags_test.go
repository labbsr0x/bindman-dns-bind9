package manager

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestBingFlags(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	err := command.ParseFlags([]string{
		"--dns-ttl=10s",
		"--dns-removal-delay=10s",
	})
	require.NoError(t, err)

	b := &Options{}
	b.InitFromViper(v)

	assert.Equal(t, time.Second*10, b.TTL)
	assert.Equal(t, time.Second*10, b.RemovalDelay)
}

func TestDefaultValues(t *testing.T) {
	v := viper.New()
	command := cobra.Command{}
	AddFlags(command.Flags())
	_ = v.BindPFlags(command.Flags())

	b := &Options{}
	b.InitFromViper(v)

	assert.Equal(t, defaultDnsTtl, b.TTL)
	assert.Equal(t, defaultDnsRemovalDelay, b.RemovalDelay)
}
