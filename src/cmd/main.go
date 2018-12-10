package main

import (
	"github.com/labbsr0x/bindman-dns-bind9/src/manager"
	"github.com/labbsr0x/bindman-dns-webhook/src/hook"
)

func main() {
	hook.Initialize(manager.New())
}
