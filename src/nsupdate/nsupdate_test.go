package nsupdate

import (
	"fmt"
	"github.com/labbsr0x/bindman-dns-webhook/src/types"
	"testing"
)

func TestNew(t *testing.T) {

}

func TestRemoveRR(t *testing.T) {

}

func TestAddRR(t *testing.T) {

}

func TestExecuteCommand(t *testing.T) {

}

func TestBuildCmdFile(t *testing.T) {

}

func TestExecCmdFile(t *testing.T) {

}

func TestCheck(t *testing.T) {
	errMsg := `The "%v" must be specified`
	errorMsgNsAddress := fmt.Sprintf(errMsg, "nameserver address")
	errorMsgKeyFileName := fmt.Sprintf(errMsg, "nameserver key file name")
	errorMsgDnsZone := fmt.Sprintf(errMsg, "DNS zone")

	type returnValue struct {
		success bool
		errs    []string
	}
	testCases := []struct {
		name     string
		nsUpdate NSUpdate
		expected returnValue
	}{
		{
			"all OK",
			NSUpdate{&Builder{Server: "localhost", KeyFile: "Ktest.com.+157+50086.key", Zone: "test.com"}},
			returnValue{true, []string{}},
		},
		{
			"all required fields",
			NSUpdate{&Builder{}},
			returnValue{false, []string{errorMsgNsAddress, errorMsgKeyFileName, errorMsgDnsZone}},
		},
		{
			"nameserver address required",
			NSUpdate{&Builder{KeyFile: "Ktest.com.+157+50086.key", Zone: "test.com"}},
			returnValue{false, []string{errorMsgNsAddress}},
		},
		{
			"nameserver key file name required",
			NSUpdate{&Builder{Server: "localhost", Zone: "test.com"}},
			returnValue{false, []string{errorMsgKeyFileName}},
		},
		{
			"",
			NSUpdate{&Builder{Server: "localhost", KeyFile: "com.+157+50086", Zone: "test.com"}},
			returnValue{false, []string{fmt.Sprintf("nameserver key file name did not match the regex %v: %s", keyFileNamePattern, "com.+157+50086")}},
		},
		{
			"DNS zone required",
			NSUpdate{&Builder{Server: "localhost", KeyFile: "Ktest.com.+157+50086.key"}},
			returnValue{false, []string{errorMsgDnsZone}},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			success, errs := test.nsUpdate.check()
			if success != test.expected.success {
				t.Errorf("It was expected success=false but returned true")
			}
			if len(errs) != len(test.expected.errs) {
				t.Errorf("The error array length must be %d but got %d", len(test.expected.errs), len(errs))
				t.FailNow()
			}
			for i, err := range test.expected.errs {
				if errs[i] != err {
					t.Errorf("Expected message was %s but got %s", err, errs[i])
				}
			}
		})
	}
}

func TestGetKeyFilePath(t *testing.T) {
	nsUpdate := NSUpdate{&Builder{Zone: "test.com", Server: "localhost", Debug: true, Port: "53", KeyFile: "Ktest.com.+157+50086.key"}}
	nsUpdate.BasePath = "data"

	expected := nsUpdate.BasePath + "/" + nsUpdate.KeyFile
	actual := nsUpdate.getKeyFilePath()
	if expected != nsUpdate.getKeyFilePath() {
		t.Errorf("expected %v but got %v", expected, actual)
	}

	nsUpdate.BasePath = "./data"
	actual = nsUpdate.getKeyFilePath()
	if expected != nsUpdate.getKeyFilePath() {
		t.Errorf("expected %v but got %v", expected, actual)
	}
}

func TestGetSubdomainName(t *testing.T) {
	nsUpdate := NSUpdate{&Builder{Zone: "test.com."}}

	tests := []struct {
		domain   string
		expected string
	}{
		{".test.com.", ""},
		{"a.test.com.", "a"},
		{"www.test.com.", "www"},
		{"www1.test.com.", "www1"},
		{"test.com.", "test.com."},
		{"_.test.com.", "_"},
		{"subdomain.test.com.", "subdomain"},
		{"subdomain.test.com.br.", "subdomain.test.com.br."},
		{"subdomain.subdomain.test.com.", "subdomain.subdomain"},
		{"subdomain.teste.com", "subdomain.teste.com"},
		{"subdomain.teste.com.", "subdomain.teste.com."},
		{"subdomain.etest.com", "subdomain.etest.com"},
		{"subdomain.etest.com.", "subdomain.etest.com."},
		{"subdomain.teste.com.br.", "subdomain.teste.com.br."},
	}

	for _, test := range tests {
		t.Run(test.domain, func(t *testing.T) {
			if subdomain := nsUpdate.getSubdomainName(test.domain); test.expected != subdomain {
				t.Errorf(`the value returned to %v was "%v" but must be "%v"`, test.domain, subdomain, test.expected)
			}
		})
	}
}

func TestCheckName(t *testing.T) {
	nsUpdate := NSUpdate{&Builder{Zone: "test.com.", Server: "localhost", Debug: true, Port: "53", KeyFile: "Ktest.com.+157+50086.key"}}

	errorMsg := "the record name '%s' is not allowed. Must obey the following pattern: '<subdomain>.%s'"

	tests := []struct {
		name     string
		expected error
	}{
		{"teste.io.", types.BadRequestError(fmt.Sprintf(errorMsg, "teste.io.", nsUpdate.Zone), nil)},
		{".test.com", types.BadRequestError(fmt.Sprintf(errorMsg, ".test.com", nsUpdate.Zone), nil)},
		{"test.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "test.com.", nsUpdate.Zone), nil)},
		{"subdomain.test.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.test.com", nsUpdate.Zone), nil)},
		{"subdomain.test.com.br", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.test.com.br", nsUpdate.Zone), nil)},
		{"subdomain.subdomain.test.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.subdomain.test.com", nsUpdate.Zone), nil)},
		{"subdomain.teste.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com", nsUpdate.Zone), nil)},
		{"subdomain.teste.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com.", nsUpdate.Zone), nil)},
		{"subdomain.etest.com", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.etest.com", nsUpdate.Zone), nil)},
		{"subdomain.etest.com.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.etest.com.", nsUpdate.Zone), nil)},
		{"subdomain.teste.com.br.", types.BadRequestError(fmt.Sprintf(errorMsg, "subdomain.teste.com.br.", nsUpdate.Zone), nil)},
		{"subdomain.subdomain.test.com.", nil},
		{"subdomain.test.com.", nil},
		{"a.test.com.", nil},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := nsUpdate.checkName(test.name)
			if test.expected == nil {
				if err != nil {
					t.Errorf("got = %v, want %v", err, test.expected)
				}
			} else {
				if err.Error() != test.expected.Error() {
					t.Errorf("got = %v, want %v", err, test.expected)
				}
			}
		})
	}

}
