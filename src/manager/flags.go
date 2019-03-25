package manager

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

const (
	dnsTtl                 = "dns-ttl"
	dnsRemovalDelay        = "dns-removal-delay"
	defaultDnsTtl          = time.Hour
	defaultDnsRemovalDelay = 10 * time.Minute
)

// AddFlags adds flags for Options.
func AddFlags(flags *pflag.FlagSet) {
	flags.Duration(dnsTtl, defaultDnsTtl, "DNS recording rule expiration time (or time-to-live)")
	flags.Duration(dnsRemovalDelay, defaultDnsRemovalDelay, "Delay in minutes to be applied to the removal of an DNS entry. This is to guarantee that in fact the removal should be processed.")
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Builder) InitFromViper(v *viper.Viper) *Builder {
	b.TTL = v.GetDuration(dnsTtl)
	b.RemovalDelay = v.GetDuration(dnsRemovalDelay)
	return b
}
