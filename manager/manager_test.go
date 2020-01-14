package manager

import (
	"fmt"
	"os"
	"sync/atomic"
	"testing"
	"time"

	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
)

const basePath = "./data"

func TestMain(m *testing.M) {
	exitCode := m.Run()
	os.Remove(basePath)
	os.Exit(exitCode)
}

func TestNew(t *testing.T) {
	if _, err := new(Builder).New(nil, ""); err == nil || err.Error() != "not possible to start the Bind9Manager; Bind9Manager expects a valid non-nil DNSUpdater" {
		t.Error("builder.New should return error with message 'not possible to start the Bind9Manager; Bind9Manager expects a valid non-nil DNSUpdater' in face of a non-valid DNSUpdater")
	}

	if _, err := new(Builder).New(new(MockDNSUpdater), ""); err == nil || err.Error() != "not possible to start the Bind9Manager; Bind9Manager expects a non-empty basePath" {
		t.Error("builder.New should return error with message 'not possible to start the Bind9Manager; Bind9Manager expects a non-empty basePath' in face of a non-valid basePath")
	}

	manager, err := new(Builder).New(new(MockDNSUpdater), basePath)
	if err != nil {
		t.Error("manager.New should not return an error in face of a valid DNSUpdater and a valid basePath")
	}
	if manager == nil {
		t.Error("manager.New should return a non-nil Bind9Manager in face of a valid DNSUpdater and a valid basePath")
	}
}

func TestAddDNSRecordAndGetAndList(t *testing.T) {
	m, _, rs := initManagerWithNRecords(1, t)

	// test get
	recordGet, err := m.GetDNSRecord(rs[0].Name, rs[0].Type)
	if recordGet == nil || err != nil {
		t.Errorf("Expecting the get of the record '%v' and type '%v' to succeed. Got result '%v' and err '%v'", rs[0].Name, rs[0].Type, recordGet, err)
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

func TestUpdateDNSRecord(t *testing.T) {
	newValue := "127.0.0.1"

	m, _, rs := initManagerWithNRecords(1, t)
	record := rs[0]
	record.Value = newValue

	err := m.UpdateDNSRecord(record)
	if err != nil {
		t.Errorf("Expecting the update of the record '%v' to succeed. Got err '%v'", record, err)
	}

	// test get
	recordGet, err := m.GetDNSRecord(record.Name, record.Type)
	if recordGet == nil || err != nil {
		t.Errorf("Expecting the get of the record '%v' and type '%v' to succeed. Got result '%v' and err '%v'", record.Name, record.Type, recordGet, err)
	}

	if recordGet != nil && (recordGet.Name != record.Name || recordGet.Value != record.Value || recordGet.Type != record.Type) {
		t.Errorf("Expecting the record saved to be equal to the the one added. Got '%v'", recordGet)
	}

	// test list
	records, err := m.GetDNSRecords()
	if records == nil || err != nil {
		t.Errorf("Expecting the list of records to return successfully. Got result nil and err '%v'", err)
	}

	if len(records) != 1 {
		t.Errorf("Expecting the list of records to have exactly one entry. Got %v", len(records))
	}

	recordRetrievedFromList := records[0]
	if recordRetrievedFromList.Name != record.Name || recordRetrievedFromList.Value != record.Value || recordRetrievedFromList.Type != record.Type {
		t.Errorf("Expecting the array entry to be exactly the same as the one updated. Got %v", recordRetrievedFromList)
	}

}

func TestDelayRemove(t *testing.T) {
	m, updater, rs := initManagerWithNRecords(1, t)

	m.RemovalDelay = 2 * time.Second
	// rest remove
	err := m.RemoveDNSRecord("test0.test.com", "A")
	if err != nil {
		t.Errorf("Expecting removal of the record '%v' to succeed. Got err '%v'", "test0.test.com", err)
	}

	if atomic.LoadUint64(&updater.RemovalCount) != 0 {
		t.Errorf("Expecting the updater.RemoveRR to not be called while IN the grace period. Got '%v' calls instead", updater.RemovalCount)
	}

	time.Sleep(3 * time.Second)

	r, err := m.GetDNSRecord(rs[0].Name, rs[0].Type)
	if err == nil {
		t.Fatalf("Must return an error when trying to get a nonexistent record, got %v", err)
	}
	if r != nil {
		t.Errorf("Expecting the removal of the file record to succeed. Got the record '%v' instead.", r)
	}

	if atomic.LoadUint64(&updater.RemovalCount) != 1 {
		t.Errorf("Expecting the updater.RemoveRR to be called exactly once after the grace period. Got '%v'", updater.RemovalCount)
	}

	// remove nonexistent record
	err = m.RemoveDNSRecord("test0.test.com", "A")
	if err == nil {
		t.Errorf("Expecting removal of the record '%v' to fail. Got err nil", "test0.test.com")
	}
	if atomic.LoadUint64(&updater.RemovalCount) != 1 {
		t.Errorf("Expecting the updater.RemoveRR to be called exactly twice after the grace period. Got '%v'", updater.RemovalCount)
	}
}

func TestGetRecordFileName(t *testing.T) {
	m, _, _ := initManagerWithNRecords(0, t)

	f := m.getRecordFileName("teste", "A")
	e := fmt.Sprintf("teste.A." + Extension)
	if f != e {
		t.Errorf("Expecting teste%v; got %v", e, f)
	}
}

func TestGetRecordName(t *testing.T) {
	m, _, _ := initManagerWithNRecords(0, t)

	name, recordType := m.getRecordNameAndType("teste.A.bindman")
	e := "teste"
	if name != e {
		t.Errorf("Expecting %v; got %v", e, name)
	}
	e = "A"
	if recordType != e {
		t.Errorf("Expecting %v; got %v", e, recordType)
	}
}

func initManagerWithNRecords(numberOfRecords int, t *testing.T) (*Bind9Manager, *MockDNSUpdater, []hookTypes.DNSRecord) {
	updater := new(MockDNSUpdater)
	updater.Result = true
	m, _ := new(Builder).New(updater, basePath)
	records := make([]hookTypes.DNSRecord, 0)

	for i := 0; i < numberOfRecords; i++ {
		record2Add := hookTypes.DNSRecord{Name: fmt.Sprintf("test%v.test.com", i), Value: "0.0.0.0", Type: "A"}
		err := m.AddDNSRecord(record2Add)
		if err != nil {
			t.Errorf("Expecting the addition of the record '%v' to succeed. Got err '%v'", record2Add, err)
		}
		records = append(records, record2Add)
	}

	return m, updater, records
}

// MockDNSUpdater defines a mock NSUpdate for unit testing the manager
type MockDNSUpdater struct {
	Result       bool
	Error        error
	RemovalCount uint64
}

func (mnsu *MockDNSUpdater) AddRR(record hookTypes.DNSRecord, ttl time.Duration) error {
	return mnsu.Error
}

func (mnsu *MockDNSUpdater) RemoveRR(name, recordType string) error {
	atomic.AddUint64(&mnsu.RemovalCount, 1)
	return mnsu.Error
}

func (mnsu *MockDNSUpdater) UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) error {
	return mnsu.Error
}
