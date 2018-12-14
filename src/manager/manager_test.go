package manager

import (
	"testing"
	"time"
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
		manager := New(new(MockNSUpdate), "./data")

		if manager.TTL != 3600 {
			t.Errorf("Default manager.TTL should be 3600 not %v", manager.TTL)
		}

		if manager.RemovalDelay != time.Duration(10*time.Minute) {
			t.Errorf("Default removal delay should be 10 min not %v", manager.RemovalDelay)
		}
	}()
}

func TestGetDNSRecord(t *testing.T) {

}

func TestAddDNSRecord(t *testing.T) {

}

func TestRemoveDNSRecord(t *testing.T) {

}

func TestDelayRemove(t *testing.T) {

}

func TestRemoveRecord(t *testing.T) {

}

func TestGetRecordFileName(t *testing.T) {

}

func TestGetRecordName(t *testing.T) {

}

// MockNSUpdate defines a mock NSUpdate for unit testing the manager
type MockNSUpdate struct {
	Result bool
	Error  error
}

func (mnsu *MockNSUpdate) RemoveRR(name, recordType string) (bool, error) {
	return mnsu.Result, mnsu.Error
}

func (mnsu *MockNSUpdate) AddRR(name, recordType, value string, ttl int) (bool, error) {
	return mnsu.Result, mnsu.Error
}
