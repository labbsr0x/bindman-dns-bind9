package manager

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/labbsr0x/bindman-dns-bind9/src/nsupdate"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/peterbourgon/diskv"
)

// Bind9Manager holds the information for managing a bind9 dns server
type Bind9Manager struct {
	DNSRecords   *diskv.Diskv
	TTL          int
	RemovalDelay time.Duration
	NSUpdate     *nsupdate.NSUpdate
}

// New creates a new Bind9Manager
func New() (result *Bind9Manager) {
	nsu, err := nsupdate.New(BasePath)
	hookTypes.PanicIfError(hookTypes.Error{Message: "Not possible to start the Bind9Manager; something went wrong while setting up NSUpdate: %s", Code: ErrInitNSUpdate, Err: err})

	result = &Bind9Manager{
		DNSRecords: diskv.New(diskv.Options{
			BasePath:     BasePath,
			Transform:    func(s string) []string { return []string{} },
			CacheSizeMax: 1024 * 1024,
		}),
		TTL:          3600,
		RemovalDelay: 10 * time.Minute,
		NSUpdate:     nsu,
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
	succ, err := m.NSUpdate.AddRR(record.Name, record.Type, record.Value, m.TTL)
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
	logrus.Infof("%s scheduled to be removed in %v seconds", name, m.RemovalDelay)
	return true, nil
}
