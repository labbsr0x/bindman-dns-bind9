package nsupdate

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	nameServerPrefix      = "nameserver."
	nameServerAddress     = nameServerPrefix + "address"
	nameServerPort        = nameServerPrefix + "port"
	nameServerKeyFile     = nameServerPrefix + "key-file"
	nameServerZone        = nameServerPrefix + "zone"
	debug                 = "debug"
	defaultNameServerPort = 53
)

// AddFlags adds flags for Options.
func AddFlags(flags *pflag.FlagSet) {
	flags.String(nameServerAddress, "", "Address of the nameserver that an instance of a Bindman will manage")
	flags.Int(nameServerPort, defaultNameServerPort, "Custom port for communication with the nameserver")
	flags.String(nameServerKeyFile, "", `Zone key-file name that will be used to authenticate with the nameserver. MUST match the regexp "K.*\.\+157\+.*\.key" and MUST be inside the /data volume`)
	flags.String(nameServerZone, "", "The name of the zone a bindman-dns-bind9 instance is able to manage")
	flags.BoolP(debug, "d", false, "The name of the zone a bindman-dns-bind9 instance is able to manage")
}

// InitFromViper initializes Options with properties retrieved from Viper.
func (b *Options) InitFromViper(v *viper.Viper) *Options {
	b.Server = v.GetString(nameServerAddress)
	b.Port = v.GetString(nameServerPort)
	b.KeyFile = v.GetString(nameServerKeyFile)
	b.Zone = v.GetString(nameServerZone)
	b.Debug = v.GetBool(debug)
	return b
}
