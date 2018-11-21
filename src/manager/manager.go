package manager

import (
	"encoding/json"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labbsr0x/sandman-bind9-manager/src/nsupdate"
	hookTypes "github.com/labbsr0x/sandman-dns-webhook/src/types"
	"github.com/peterbourgon/diskv"
)

// Bind9Manager holds the information for managing a bind9 dns server
type Bind9Manager struct {
	DNSRecords         *diskv.Diskv
	TTL                int
	RemovalWaitingTime time.Duration
	NSUpdate           *nsupdate.NSUpdate
}

// New creates a new Bind9Manager
func New() *Bind9Manager {
	nsu, err := nsupdate.New("", "", "")
	hookTypes.PanicIfError(hookTypes.Error{Message: "Not possible to start the Bind9Manager; something went wrong while setting up NSUpdate: %s", Code: ErrInitNSUpdate, Err: err})
	return &Bind9Manager{
		DNSRecords: diskv.New(diskv.Options{
			BasePath:     "data",
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
		TTL:                3600,             // TODO: get from env
		RemovalWaitingTime: 10 * time.Second, // TODO: get from env
		NSUpdate:           nsu,              // TODO: get from env
	}
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Bind9Manager) GetDNSRecords() ([]hookTypes.DNSRecord, error) {
	toReturn := []hookTypes.DNSRecord{}
	// TODO: navigate data subfolders
	return toReturn, nil
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Bind9Manager) GetDNSRecord(name string) (*hookTypes.DNSRecord, error) {
	r, _ := m.DNSRecords.Read(name)
	toReturn := &hookTypes.DNSRecord{}
	err := json.Unmarshal(r, toReturn)
	if err == nil {
		return toReturn, nil
	}
	return nil, err
}

// AddDNSRecord adds a new DNS record
func (m *Bind9Manager) AddDNSRecord(record hookTypes.DNSRecord) (bool, error) {
	succ, err := m.NSUpdate.AddRR(record.Name, record.IPAddr, record.TTL)
	if succ {
		r, _ := json.Marshal(record)
		m.DNSRecords.Write(record.Name, r)
		return true, nil
	}
	return false, err
}

// RemoveDNSRecord removes a DNS record
func (m *Bind9Manager) RemoveDNSRecord(name string) (bool, error) {
	go m.delayRemove(name)
	logrus.Infof("%s scheduled to be removed in %v seconds", name, m.RemovalWaitingTime)
	return true, nil
}

// delayRemove schedules the removal of a DNS Resource Record
// it cancels the operation when it idenfities the name was readded
func (m *Bind9Manager) delayRemove(name string) {
	m.DNSRecords.Erase(name) // marks its removal
	c := time.Tick(m.RemovalWaitingTime)
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
