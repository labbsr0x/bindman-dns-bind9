package nsupdate

import (
	"fmt"
	"path"
	"regexp"
	"strings"
)

const (
	// SANDMAN_NAMESERVER_ADDRESS environment variable identifier for the nameserver address
	SANDMAN_NAMESERVER_ADDRESS = "SANDMAN_NAMESERVER_ADDRESS"

	// SANDMAN_NAMESERVER_PORT environment variable identifier for the nameserver port
	SANDMAN_NAMESERVER_PORT = "SANDMAN_NAMESERVER_PORT"

	// SANDMAN_NAMESERVER_KEYFILE environment variable identifier for the nameserver key name
	SANDMAN_NAMESERVER_KEYFILE = "SANDMAN_NAMESERVER_KEYFILE"

	// SANDMAN_NAMESERVER_ZONE environment variable identifier for the zone to be managed
	SANDMAN_NAMESERVER_ZONE = "SANDMAN_NAMESERVER_ZONE"
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

	if nsu.Zone == "" {
		errs = append(errs, fmt.Sprintf(errMsg, SANDMAN_NAMESERVER_ZONE))
	}

	m := `K.*\.\+157\.\+.*\.key`
	if succ, _ := regexp.MatchString(nsu.KeyFile, m); !succ {
		errs = append(errs, "Environment variable %s did not match the regex %v: %s", SANDMAN_NAMESERVER_KEYFILE, m)
	}

	// TODO: Test connection

	return false, errs
}

// getKeyFilePath joins the base path with key file name
func (nsu *NSUpdate) getKeyFilePath() string {
	return path.Join(nsu.BasePath, nsu.KeyFile)
}

// getSubdomainName we expect names to come in the format subdomain.zone. This function returns the subdomain part
func (nsu *NSUpdate) getSubdomainName(name string) string {
	str := strings.Replace(name, nsu.Zone, "", 1)
	return strings.TrimSuffix(str, ".")
}

// checkName checks if the name is in the expected format: subdomain.zone
func (nsu *NSUpdate) checkName(name string) (bool, error) {
	m := fmt.Sprintf(".*\\.%s", nsu.Zone)
	return regexp.MatchString(name, m)
}
