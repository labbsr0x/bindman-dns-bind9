package nsupdate

import (
	"bufio"
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
	RemoveRR(name, recordType string) (err error)
	AddRR(record hookTypes.DNSRecord, ttl time.Duration) (err error)
	UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (err error)
}

// New constructs a new NSUpdate instance from environment variables
func (b *Builder) New(basePath string) (result *NSUpdate, err error) {
	b.BasePath = basePath

	if !strings.HasSuffix(b.Zone, ".") {
		b.Zone = fmt.Sprintf("%s.", b.Zone)
	}

	result = &NSUpdate{b}

	if succ, errs := result.check(); !succ {
		err = fmt.Errorf("Errors encountered:\n\t%v", strings.Join(errs, "\n\t"))
	}
	return
}

// RemoveRR removes a Resource Record
func (nsu *NSUpdate) RemoveRR(name, recordType string) (err error) {
	err = nsu.checkName(name)
	if err == nil {
		cmd := nsu.buildDeleteCommand(name, recordType)
		logrus.Infof("cmd to be executed: %s", cmd)
		err = nsu.ExecuteCommand(cmd)
	}
	return
}

// AddRR adds a Resource Record
func (nsu *NSUpdate) AddRR(record hookTypes.DNSRecord, ttl time.Duration) (err error) {
	err = nsu.checkName(record.Name)
	if err == nil {
		cmd := nsu.buildAddCommand(record.Name, record.Type, record.Value, ttl)
		logrus.Infof("cmd to be executed: %s", cmd)
		err = nsu.ExecuteCommand(cmd)
	}
	return
}

// UpdateRR updates a DNS Resource Record
func (nsu *NSUpdate) UpdateRR(record hookTypes.DNSRecord, ttl time.Duration) (err error) {
	err = nsu.checkName(record.Name)
	if err == nil {
		deleteCmd := nsu.buildDeleteCommand(record.Name, record.Type)
		addCmd := nsu.buildAddCommand(record.Name, record.Type, record.Value, ttl)
		cmd := fmt.Sprintf("%v\n%v", deleteCmd, addCmd)
		logrus.Infof("cmd to be executed: %s", cmd)
		err = nsu.ExecuteCommand(cmd)
	}
	return
}

// ExecuteCommand executes a given nsupdate command
func (nsu *NSUpdate) ExecuteCommand(cmd string) (err error) {
	fileName, err := nsu.BuildCmdFile(cmd)
	if err == nil {
		logrus.Infof("Created the nsupdate cmd file %s successfully", fileName)
		err = nsu.ExecCmdFile(fileName)
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

		_, err = writer.WriteString(fmt.Sprintf("server %s\n", nsu.Server))
		_, err = writer.WriteString(fmt.Sprintf("zone %s\n", nsu.Zone))
		_, err = writer.WriteString(cmd + "\n")
		_, err = writer.WriteString("send")

		err = writer.Flush()
		err = f.Close()

		fileName = f.Name()
	}
	return
}

// ExecCmdFile executes an nsupdate cmd file
func (nsu *NSUpdate) ExecCmdFile(filePath string) (err error) {
	keyFilePath := nsu.getKeyFilePath()
	// The -v option makes nsupdate use a TCP connection.
	exe := exec.Command("nsupdate", "-v", "-k", keyFilePath, filePath)
	msg, err := exe.CombinedOutput()

	if err != nil {
		err = fmt.Errorf("error executing command file %s: %s %s", exe.Path, err.Error(), string(msg))
	}
	return
}
