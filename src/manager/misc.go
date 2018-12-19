package manager

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/sirupsen/logrus"
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
	Extension = "bindman"
)

// delayRemove schedules the removal of a DNS Resource Record
// it cancels the operation when it idenfities the name was readded
func (m *Bind9Manager) delayRemove(name, recordType string) {
	record, err := m.GetDNSRecord(name, recordType)
	if err == nil {
		go m.removeRecord(name, recordType) // marks its removal intent
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

// saveRecord saves a record to the local storage
func (m *Bind9Manager) saveRecord(record hookTypes.DNSRecord) (err error) {
	var r []byte
	r, err = json.Marshal(record)
	if err == nil {
		m.Door.Lock()
		defer m.Door.Unlock()

		err = m.DNSRecords.Write(m.getRecordFileName(record.Name, record.Type), r)
	}
	return
}

// removeRecord removes the record
func (m *Bind9Manager) removeRecord(recordName, recordType string) {
	m.Door.Lock()
	defer m.Door.Unlock()
	m.DNSRecords.Erase(m.getRecordFileName(recordName, recordType)) // marks its removal
}

// getRecordFileName return the name of the file holding the record information
func (m *Bind9Manager) getRecordFileName(recordName, recordType string) string {
	toReturn := fmt.Sprintf("%v.%v.%v", recordName, recordType, Extension)
	return toReturn
}

// getRecordName returns the name of a record from its fileName
func (m *Bind9Manager) getRecordNameAndType(fileName string) (string, string) {
	subName := strings.TrimSuffix(fileName, "."+Extension)
	i := strings.LastIndex(subName, ".")
	return subName[:i], subName[i+1:]
}
