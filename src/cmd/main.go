package main

import (
	"github.com/labbsr0x/sandman-bind9-manager/src/manager"
	"github.com/labbsr0x/sandman-dns-webhook/src/hook"
)

func main() {
	hook.Initialize(manager.New())
}
