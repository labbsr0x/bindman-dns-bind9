package main

import (
	"os"

	"github.com/labbsr0x/bindman-dns-bind9/src/manager"
	"github.com/labbsr0x/bindman-dns-bind9/src/nsupdate"
	"github.com/labbsr0x/bindman-dns-webhook/src/hook"
	"github.com/sirupsen/logrus"
)

func main() {
	basePath := "/data"
	nsu, err := nsupdate.New(basePath)
	if err != nil {
		logrus.Errorf("An error ocurred while setting up the DNS Manager: %v", err)
		os.Exit(manager.ErrInitNSUpdate)
	}
	hook.Initialize(manager.New(nsu, basePath))
}
