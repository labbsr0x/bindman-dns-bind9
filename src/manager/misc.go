package manager

import (
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
)

const (
	// ErrInitNSUpdate error code for problems while setting up NSUpdate
	ErrInitNSUpdate = iota
)

const (
	// SANDMAN_DNS_TTL environment variable identifier for the time-to-live to be applied
	SANDMAN_DNS_TTL = "BINDMAN_DNS_TTL"

	// SANDMAN_DNS_REMOVAL_DELAY environment variable identifier for the removal delay time to be applied
	SANDMAN_DNS_REMOVAL_DELAY = "BINDMAN_DNS_REMOVAL_DELAY"

	// Extension sets the extension of the files holding the records infos
	Extension = ".bindman"
)

// delayRemove schedules the removal of a DNS Resource Record
// it cancels the operation when it idenfities the name was readded
func (m *Bind9Manager) delayRemove(name string) {
	record, err := m.GetDNSRecord(name) // marks its removal intent
	if err == nil {
		go m.removeRecord(name)
		c := time.Tick(m.RemovalDelay)
		for {
			select {
			case <-c:
				if _, err := m.DNSRecords.Read(name); err == nil { // record has been readded
					logrus.Infof("Cancelling delayed removal of '%s'", name)
					return
				}

				// only remove in case the record has not been readded
				if succ, err := m.DNSUpdater.RemoveRR(name, record.Type); !succ {
					logrus.Infof("Error occurred while trying to remove '%s': %s", name, err)
				}
				return
			}
		}
	} else {
		logrus.Errorf("Service '%v' cannot be removed given it does not exist.", name)
	}
}

// removeRecord removes the record
func (m *Bind9Manager) removeRecord(name string) {
	m.Door.Lock()
	defer m.Door.Unlock()
	m.DNSRecords.Erase(m.getRecordFileName(name)) // marks its removal
}

// getRecordFileName return the name of the file holding the record information
func (m *Bind9Manager) getRecordFileName(recordName string) string {
	return recordName + Extension
}

func (m *Bind9Manager) getRecordName(fileName string) string {
	return strings.TrimSuffix(fileName, Extension)
}
