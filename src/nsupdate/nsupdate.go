package nsupdate

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
)

// NSUpdate holds the information necessary to successfully run nsupdate requests
type NSUpdate struct {
	Server   string
	Port     string
	KeyFile  string
	BasePath string
	Zone     string
}

// New constructs a new NSUpdate instance from environment variables
func New(basePath string) (result *NSUpdate, err error) {
	result = &NSUpdate{
		Server:   strings.Trim(os.Getenv(SANDMAN_NAMESERVER_ADDRESS), " "),
		Port:     strings.Trim(os.Getenv(SANDMAN_NAMESERVER_PORT), " "),
		KeyFile:  strings.Trim(os.Getenv(SANDMAN_NAMESERVER_KEYFILE), " "),
		Zone:     strings.Trim(os.Getenv(SANDMAN_NAMESERVER_ZONE), " "),
		BasePath: basePath,
	}

	if succ, errs := result.check(); !succ {
		err = fmt.Errorf("Errors encountered: %v", strings.Join(errs, ", "))
	}

	return
}

// RemoveRR removes a Resource Record
func (nsu *NSUpdate) RemoveRR(name string) (bool, error) {
	cmd := fmt.Sprintf("update delete %s.%s. %s\n", nsu.getSubdomainName(name), nsu.Zone, "A")
	logrus.Infof("cmd to be executed: %s", cmd)
	return nsu.ExecuteCommand(cmd)
}

// AddRR adds a Resource Record
func (nsu *NSUpdate) AddRR(name string, ipaddr string, ttl int) (success bool, err error) {
	success, err = nsu.checkName(name)
	if success {
		cmd := fmt.Sprintf("update add %s.%s. %d %s %s\n", nsu.getSubdomainName(name), nsu.Zone, ttl, "A", ipaddr)
		logrus.Infof("cmd to be executed: %s", cmd)
		success, err = nsu.ExecuteCommand(cmd)
	} else {
		err = fmt.Errorf("The record name '%s' is not allowed. Must be of the following format: <subdomain>.%s", name, nsu.Zone)
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
		os.Remove(fileName)
	}
	return
}

// BuildCmdFile creates an nsupdate cmd file
func (nsu *NSUpdate) BuildCmdFile(cmd string) (fileName string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), "sandman"+uuid.New().String())
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
