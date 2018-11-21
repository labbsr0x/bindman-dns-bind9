package nsupdate

import (
	"fmt"
	"os"
	"strings"
)

// NSUpdate holds the information necessary to successfully run nsupdate requests
type NSUpdate struct {
	Server   string
	Port     string
	KeyFile  string
	BasePath string
}

// New constructs a new NSUpdate instance from environment variables
func New(basePath string) (result *NSUpdate, err error) {
	result = &NSUpdate{
		Server:   strings.Trim(os.Getenv(SANDMAN_NAMESERVER_ADDRESS), " "),
		Port:     strings.Trim(os.Getenv(SANDMAN_NAMESERVER_PORT), " "),
		KeyFile:  strings.Trim(os.Getenv(SANDMAN_NAMESERVER_KEYFILE), " "),
		BasePath: basePath,
	}

	if succ, errs := result.check(); !succ {
		err = fmt.Errorf("Errors encountered: %v", strings.Join(errs, ", "))
	}

	return
}

// RemoveRR removes a Resource Record
func (nsu *NSUpdate) RemoveRR(name string) (success bool, err error) {
	success = true
	err = nil
	// TODO
	return
}

// AddRR adds a Resource Record
func (nsu *NSUpdate) AddRR(name string, ipaddr string, ttl int) (success bool, err error) {
	success = true
	err = nil
	// TODO
	return
}
