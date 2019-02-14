package nsupdate

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	hookTypes "github.com/labbsr0x/bindman-dns-webhook/src/types"
	"github.com/sirupsen/logrus"
)

type Builder struct {
	Server   string
	Port     string
	KeyFile  string
	BasePath string
	Zone     string
	Debug    bool
}

// NSUpdate holds the information necessary to successfully run nsupdate requests
type NSUpdate struct {
	*Builder
}

// DNSUpdater defines an interface to communicate with DNS Server via nsupdate commands
type DNSUpdater interface {
	RemoveRR(name, recordType string) (success bool, err error)
	AddRR(name, recordType, value string, ttl time.Duration) (success bool, err error)
	UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (success bool, err error)
}

// New constructs a new NSUpdate instance from environment variables
func (b *Builder) New(basePath string) (result *NSUpdate, err error) {
	b.BasePath = basePath
	result = &NSUpdate{b}

	if succ, errs := result.check(); !succ {
		err = fmt.Errorf("Errors encountered: %v", strings.Join(errs, ", "))
	}
	return
}

// RemoveRR removes a Resource Record
func (nsu *NSUpdate) RemoveRR(name, recordType string) (succ bool, err error) {
	cmd, err := nsu.buildDeleteCommand(name, recordType)
	if err == nil {
		logrus.Infof("cmd to be executed: %s", cmd)
		succ, err = nsu.ExecuteCommand(cmd)
	}
	return
}

// AddRR adds a Resource Record
func (nsu *NSUpdate) AddRR(name, recordType, value string, ttl time.Duration) (succ bool, err error) {
	cmd, err := nsu.buildAddCommand(name, recordType, value, ttl)
	if err == nil {
		logrus.Infof("cmd to be executed: %s", cmd)
		succ, err = nsu.ExecuteCommand(cmd)
	}
	return
}

// UpdateRR updates a DNS Resource Record
func (nsu *NSUpdate) UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (succ bool, err error) {
	deleteCmd, err := nsu.buildDeleteCommand(record.Name, record.Type)
	if err == nil {
		addCmd, err := nsu.buildAddCommand(record.Name, record.Type, record.Value, ttl)
		if err == nil {
			cmd := fmt.Sprintf("%v\n%v", deleteCmd, addCmd)
			logrus.Infof("cmd to be executed: %s", cmd)
			succ, err = nsu.ExecuteCommand(cmd)
		}
	}
	return
}

// ExecuteCommand executes a given nsupdate command
func (nsu *NSUpdate) ExecuteCommand(cmd string) (success bool, err error) {
	fileName, err := nsu.BuildCmdFile(cmd)
	if err == nil {
		logrus.Infof("Created the nsupdate cmd file %s successfully", fileName)
		success, err = nsu.ExecCmdFile(fileName)
		if err == nil {
			logrus.Infof("Executes cmd %s successfully", cmd)
		}
		if !nsu.Debug {
			os.Remove(fileName)
		}
	}
	return
}

// BuildCmdFile creates an nsupdate cmd file
func (nsu *NSUpdate) BuildCmdFile(cmd string) (fileName string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), "bindman-"+uuid.New().String())
	if err == nil {
		writer := bufio.NewWriter(f)

		writer.WriteString(fmt.Sprintf("server %s\n", nsu.Server))
		writer.WriteString(fmt.Sprintf("zone %s\n", nsu.Zone))
		writer.WriteString(cmd + "\n")
		writer.WriteString("send\n")

		writer.Flush()
		f.Close()

		fileName = f.Name()
	}
	return
}

// ExecCmdFile executes an nsupdate cmd file
func (nsu *NSUpdate) ExecCmdFile(filePath string) (success bool, err error) {
	var out bytes.Buffer
	exe := exec.Command("nsupdate", "-v", "-k", nsu.getKeyFilePath(), filePath)
	exe.Stdout = &out
	err = exe.Run()

	if err == nil {
		success = true
	} else {
		err = fmt.Errorf("%s: %s", err.Error(), out.String())
	}
	return
}
