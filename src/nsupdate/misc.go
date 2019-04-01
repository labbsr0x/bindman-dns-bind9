package nsupdate

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"path"
	"regexp"
	"strings"
	"time"
)

const keyFileNamePattern = `K.*\.\+157\+.*\.key`

// check tests if a NSUpdate setup is ok; returns a set of error strings in case something is not right
func (nsu *NSUpdate) check() (success bool, errs []string) {
	errMsg := `The "%v" must be specified`

	if nsu.Server == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "nameserver address"))
	}

	if nsu.KeyFile == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "nameserver key file name"))
	} else {
		if succ, _ := regexp.MatchString(keyFileNamePattern, nsu.KeyFile); !succ {
			errs = append(errs, fmt.Sprintf("nameserver key file name did not match the regex %v: %s", keyFileNamePattern, nsu.KeyFile))
		}
	}

	if nsu.Zone == "" {
		errs = append(errs, fmt.Sprintf(errMsg, "DNS zone"))
	}

	// TODO: Test connection
	return len(errs) == 0, errs
}

// getKeyFilePath joins the base path with key file name
func (nsu *NSUpdate) getKeyFilePath() string {
	return path.Join(nsu.BasePath, nsu.KeyFile)
}

// getSubdomainName we expect names to come in the format subdomain.zone. This function returns the subdomain part
func (nsu *NSUpdate) getSubdomainName(name string) string {
	return strings.TrimSuffix(name, "."+nsu.Zone)
}

// checkName checks if the name is in the expected format: subdomain.zone
func (nsu *NSUpdate) checkName(name string) (err error) {
	if !strings.HasSuffix(name, "."+nsu.Zone) {
		err = types.BadRequestError(fmt.Sprintf("the record name '%s' is not allowed. Must obey the following pattern: '<subdomain>.%s'", name, nsu.Zone), nil)
	}
	return
}

// buildAddCommand builds a nsupdate add command
func (nsu *NSUpdate) buildAddCommand(recordName, recordType, value string, ttl time.Duration) string {
	return fmt.Sprintf("update add %s.%s %d %s %s", nsu.getSubdomainName(recordName), nsu.Zone, int(ttl.Seconds()), recordType, value)
}

// buildDeleteCommand builds a nsupdate delete command
func (nsu *NSUpdate) buildDeleteCommand(recordName, recordType string) string {
	return fmt.Sprintf("update delete %s.%s %s", nsu.getSubdomainName(recordName), nsu.Zone, recordType)
}
