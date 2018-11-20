package manager

import (
	hookTypes "github.com/labbsr0x/sandman-dns-webhook/src/types"
)

// Bind9Manager holds the information for managing a bind9 dns server
type Bind9Manager struct {
	DNSRecords map[string]hookTypes.DNSRecord
}

// New creates a new Bind9Manager
func New() *Bind9Manager {
	return &Bind9Manager{DNSRecords: make(map[string]hookTypes.DNSRecord)}
}

// GetDNSRecords retrieves all the dns records being managed
func (m *Bind9Manager) GetDNSRecords() ([]hookTypes.DNSRecord, error) {
	toReturn := []hookTypes.DNSRecord{}
	for _, v := range m.DNSRecords {
		toReturn = append(toReturn, v)
	}
	return toReturn, nil
}

// GetDNSRecord retrieves the dns record identified by name
func (m *Bind9Manager) GetDNSRecord(name string) (*hookTypes.DNSRecord, error) {
	if record, ok := m.DNSRecords[name]; ok {
		return &record, nil
	}
	return nil, nil
}

// AddDNSRecord adds a new DNS record
func (m *Bind9Manager) AddDNSRecord(record hookTypes.DNSRecord) (bool, error) {
	m.DNSRecords[record.Name] = record
	return true, nil
}

// RemoveDNSRecord removes a DNS record
func (m *Bind9Manager) RemoveDNSRecord(name string) (bool, error) {
	delete(m.DNSRecords, name)
	return true, nil
}
