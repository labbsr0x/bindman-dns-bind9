package manager

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/labbsr0x/bindman-dns-bind9/src/nsupdate"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/peterbourgon/diskv"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	TTL          time.Duration
	RemovalDelay time.Duration
}

// Bind9Manager holds the information for managing a bind9 dns server
type Bind9Manager struct {
	*Builder
	DNSRecords *diskv.Diskv
	Door       *sync.RWMutex
	DNSUpdater nsupdate.DNSUpdater
}

// New creates a new Bind9Manager
func (b *Builder) New(dnsupdater nsupdate.DNSUpdater, basePath string) (*Bind9Manager, error) {
	if dnsupdater == nil {
		return nil, errors.New("not possible to start the Bind9Manager; Bind9Manager expects a valid non-nil DNSUpdater")
	}

	if strings.TrimSpace(basePath) == "" {
		return nil, errors.New("not possible to start the Bind9Manager; Bind9Manager expects a non-empty basePath")
	}

	result := &Bind9Manager{
		DNSRecords: diskv.New(diskv.Options{
			BasePath:     basePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
		Builder:    b,
		Door:       new(sync.RWMutex),
		DNSUpdater: dnsupdater,
	}
	return result, nil
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Bind9Manager) GetDNSRecords() (records []hookTypes.DNSRecord, err error) {
	m.Door.RLock()
	defer m.Door.RUnlock()

	err = filepath.Walk(m.DNSRecords.BasePath, func(path string, info os.FileInfo, errr error) error {
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
