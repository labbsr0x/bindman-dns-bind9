package manager

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
)

func TestNew(t *testing.T) {
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("manager.New should not accept a nil DNSUpdater")
			}
		}()
		New(nil, "")
	}()

	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("manager.New should not panic in face of a valid DNSUpdater and a valid basePath")
			}
		}()
		manager := New(new(MockDNSUpdater), "./data")

		if manager.TTL != 3600 {
			t.Errorf("Default manager.TTL should be 3600 not %v", manager.TTL)
		}

		if manager.RemovalDelay != time.Duration(10*time.Minute) {
			t.Errorf("Default removal delay should be 10 min not %v", manager.RemovalDelay)
		}
	}()
}

func TestAddDNSRecordAndGetAndList(t *testing.T) {
	m, _, rs := initManagerWithNRecords(1, t)

	// test get
	recordGet, err := m.GetDNSRecord(rs[0].Name)
	if recordGet == nil || err != nil {
		t.Errorf("Expecting the get of the record '%v' to succeed. Got result '%v' and err '%v'", rs[0].Name, recordGet, err)
	}

	if recordGet != nil && (recordGet.Name != rs[0].Name || recordGet.Value != rs[0].Value || recordGet.Type != rs[0].Type) {
		t.Errorf("Expecting the record saved to be equal to the the one added. Got '%v'", recordGet)
	}

	// test list
	records, err := m.GetDNSRecords()
	if records == nil || err != nil {
		t.Errorf("Expecting the list of added records to return successfully. Got result '%v' and err '%v'", records, err)
	}

	if len(records) != 1 {
		t.Errorf("Expecting the list of records to have exactly one entry. Got %v", len(records))
	}

	if records[0].Name != rs[0].Name || records[0].Value != rs[0].Value || records[0].Type != rs[0].Type {
		t.Errorf("Expecting the array entry to be exactly the same as the one added. Got %v", records[0])
	}

}

func TestDelayRemove(t *testing.T) {
	m, updater, rs := initManagerWithNRecords(1, t)

	m.RemovalDelay = 2 * time.Second
	// rest remove
	result, err := m.RemoveDNSRecord("test0.test.com")
	if !result || err != nil {
		t.Errorf("Expecting removal of the record '%v' to succeed. Got result '%v' and err '%v'", "test0.test.com", result, err)
	}

	if atomic.LoadUint64(&updater.RemovalCount) != 0 {
		t.Errorf("Expecting the updater.RemoveRR to not be called while IN the grace period. Got '%v' calls instead", updater.RemovalCount)
	}

	time.Sleep(3 * time.Second)

	r, _ := m.GetDNSRecord(rs[0].Name)
	if r != nil {
		t.Errorf("Expecting the removal of the file record to succeed. Got the record '%v' instead.", r)
	}

	if atomic.LoadUint64(&updater.RemovalCount) != 1 {
		t.Errorf("Expecting the updater.RemoveRR to be called exactly once after the grace period. Got '%v'", updater.RemovalCount)
	}
}

func TestGetRecordFileName(t *testing.T) {
	m, _, _ := initManagerWithNRecords(0, t)

	f := m.getRecordFileName("teste")
	e := fmt.Sprintf("teste" + Extension)
	if f != e {
		t.Errorf("Expecting teste%v; got %v", e, f)
	}
}

func TestGetRecordName(t *testing.T) {
	m, _, _ := initManagerWithNRecords(0, t)

	f := m.getRecordName("teste.bindman")
	e := "teste"
	if f != e {
		t.Errorf("Expecting %v; got %v", e, f)
	}
}

func initManagerWithNRecords(numberOfRecords int, t *testing.T) (*Bind9Manager, *MockDNSUpdater, []hookTypes.DNSRecord) {
	updater := new(MockDNSUpdater)
	updater.Result = true
	m := New(updater, "./data")
	records := make([]hookTypes.DNSRecord, 0)

	for i := 0; i < numberOfRecords; i++ {
		record2Add := hookTypes.DNSRecord{Name: fmt.Sprintf("test%v.test.com", i), Value: "0.0.0.0", Type: "A"}
		result, err := m.AddDNSRecord(record2Add)
		if !result || err != nil {
			t.Errorf("Expecting the addition of the record '%v' to succeed. Got result '%v' and err '%v'", record2Add, result, err)
		}
		records = append(records, record2Add)
	}

	return m, updater, records
}

func areTwoRecordsEqual(r1, r2 hookTypes.DNSRecord) bool {
	return !(r1.Name != r2.Name || r1.Value != r2.Value || r1.Type != r2.Type)
}

// MockDNSUpdater defines a mock NSUpdate for unit testing the manager
type MockDNSUpdater struct {
	Result       bool
	Error        error
	RemovalCount uint64
}

func (mnsu *MockDNSUpdater) RemoveRR(name, recordType string) (bool, error) {
	atomic.AddUint64(&mnsu.RemovalCount, 1)
	return mnsu.Result, mnsu.Error
}

func (mnsu *MockDNSUpdater) AddRR(name, recordType, value string, ttl int) (bool, error) {
	return mnsu.Result, mnsu.Error
}
