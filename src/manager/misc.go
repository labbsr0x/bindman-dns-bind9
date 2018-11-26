package manager

import (
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	// ErrInitNSUpdate error code for problems while setting up NSUpdate
	ErrInitNSUpdate = iota
)

const (
	// BasePath defines the base path for storing and reading application specific files
	BasePath = "/data"

	// SANDMAN_DNS_TTL environment variable identifier for the time-to-live to be applied
	SANDMAN_DNS_TTL = "SANDMAN_DNS_TTL"

	// SANDMAN_DNS_REMOVAL_DELAY environment variable identifier for the removal delay time to be applied
	SANDMAN_DNS_REMOVAL_DELAY = "SANDMAN_DNS_REMOVAL_DELAY"
)

// delayRemove schedules the removal of a DNS Resource Record
// it cancels the operation when it idenfities the name was readded
func (m *Bind9Manager) delayRemove(name string) {
	m.DNSRecords.Erase(name) // marks its removal
	c := time.Tick(m.RemovalDelay)
	for {
		select {
		case <-c:
			if _, err := m.DNSRecords.Read(name); err == nil { // record has been readded
				logrus.Infof("Cancelling delayed removal of '%s'", name)
				return
			}

			// only remove in case the record has not been readded
			if succ, err := m.NSUpdate.RemoveRR(name); !succ {
				logrus.Infof("Error occurred while trying to remove '%s': %s", name, err)
			}
			return
		}
	}
}
