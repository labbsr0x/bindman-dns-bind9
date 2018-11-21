package nsupdate

import (
	"fmt"
	"regexp"
)

const (
	// SANDMAN_NAMESERVER_ADDRESS environment variable identifier for the nameserver address
	SANDMAN_NAMESERVER_ADDRESS = "SANDMAN_NAMESERVER_ADDRESS"

	// SANDMAN_NAMESERVER_PORT environment variable identifier for the nameserver port
	SANDMAN_NAMESERVER_PORT = "SANDMAN_NAMESERVER_PORT"

	// SANDMAN_NAMESERVER_KEYFILE environment variable identifier for the nameserver key name
	SANDMAN_NAMESERVER_KEYFILE = "SANDMAN_NAMESERVER_KEYFILE"
)

// check tests if a NSUpdate setup is ok; returns a set of error strings in case something is not right
func (nsu *NSUpdate) check() (success bool, errs []string) {
	errMsg := "The environment variable %s cannot be empty"
	if nsu.Server == "" {
		errs = append(errs, fmt.Sprintf(errMsg, SANDMAN_NAMESERVER_ADDRESS))
	}

	if nsu.KeyFile == "" {
		errs = append(errs, fmt.Sprintf(errMsg, SANDMAN_NAMESERVER_KEYFILE))
	}

	m := `K.*\.\+157\.\+.*\.key`
	if succ, _ := regexp.MatchString(nsu.KeyFile, m); !succ {
		errs = append(errs, "Environment variable %s did not match the regex %v: %s", SANDMAN_NAMESERVER_KEYFILE, m)
	}

	// TODO: Test connection

	return false, errs
}
