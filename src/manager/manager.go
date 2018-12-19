package manager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/labbsr0x/bindman-dns-bind9/src/nsupdate"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/peterbourgon/diskv"
	"github.com/sirupsen/logrus"
)

// Bind9Manager holds the information for managing a bind9 dns server
type Bind9Manager struct {
	BasePath     string
	DNSRecords   *diskv.Diskv
	TTL          int
	RemovalDelay time.Duration
	Door         *sync.RWMutex
	DNSUpdater   nsupdate.DNSUpdater
}

// New creates a new Bind9Manager
func New(dnsupdater nsupdate.DNSUpdater, basePath string) (result *Bind9Manager) {
	if dnsupdater == nil {
		hookTypes.Panic(hookTypes.Error{Message: "Not possible to start the Bind9Manager; Bind9Manager expects a valid non-nil DNSUpdater", Code: ErrInitNSUpdate})
	}

	result = &Bind9Manager{
		DNSRecords: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
		BasePath:     basePath,
		TTL:          3600,
		RemovalDelay: 10 * time.Minute,
		Door:         new(sync.RWMutex),
		DNSUpdater:   dnsupdater,
	}

	// get ttl from env
	if ttl, err := strconv.Atoi(strings.Trim(os.Getenv(SANDMAN_DNS_TTL), " ")); err == nil {
		result.TTL = ttl
	}

	// get removal delay from env
	if r, err := strconv.Atoi(strings.Trim(os.Getenv(SANDMAN_DNS_REMOVAL_DELAY), " ")); err == nil {
		result.RemovalDelay = time.Duration(r) * time.Minute
	}

	return result
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Bind9Manager) GetDNSRecords() (records []hookTypes.DNSRecord, err error) {
	m.Door.RLock()
	defer m.Door.RUnlock()

	err = filepath.Walk(m.BasePath, func(path string, info os.FileInfo, errr error) error {
		if strings.HasSuffix(path, Extension) {
			r, _ := m.GetDNSRecord(m.getRecordNameAndType(info.Name()))
			records = append(records, *r)
		}
		return nil
	})

	return
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Bind9Manager) GetDNSRecord(name, recordType string) (record *hookTypes.DNSRecord, err error) {
	m.Door.RLock()
	defer m.Door.RUnlock()

	var r []byte
	r, err = m.DNSRecords.Read(m.getRecordFileName(name, recordType))
	if err == nil {
		err = json.Unmarshal(r, &record)
	}
	return
}

// AddDNSRecord adds a new DNS record
func (m *Bind9Manager) AddDNSRecord(record hookTypes.DNSRecord) (succ bool, err error) {
	succ, err = m.DNSUpdater.AddRR(record.Name, record.Type, record.Value, m.TTL)
	if succ {
		err = m.saveRecord(record)
		succ = err == nil
	}
	return
}

// UpdateDNSRecord updates an existing dns record
func (m *Bind9Manager) UpdateDNSRecord(record hookTypes.DNSRecord) (succ bool, err error) {
	succ, err = m.DNSUpdater.UpdateRR(record, m.TTL)
	if succ {
		err = m.saveRecord(record)
		succ = err == nil
	}
	return
}

// RemoveDNSRecord removes a DNS record
func (m *Bind9Manager) RemoveDNSRecord(name, recordType string) (bool, error) {
	go m.delayRemove(name, recordType)
	logrus.Infof("Record '%s' with type '%v' scheduled to be removed in %v seconds", name, recordType, m.RemovalDelay)
	return true, nil
}
