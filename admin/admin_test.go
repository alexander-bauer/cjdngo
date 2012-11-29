package admin

import (
	"github.com/SashaCrofter/cjdngo"
	"testing"
)

var (
	cjdns *CJDNS
)

func TestConnect(t *testing.T) {
	conf, err := cjdngo.ReadConf("/etc/cjdroute.conf")
	if err != nil || len(conf.Admin.Password) == 0 {
		//This is not related to the test.
		t.Log(err)
		t.Log("Could not read the config file. Skipping test.")
		return
	}

	cjd, err := Connect("127.0.0.1", "11234", conf.Admin.Password)
	if err != nil {
		t.Fatal(err)
	}

	if !cjd.Ping() {
		t.Fatal("Server did not respond to ping.")
	}
	cjdns = cjd
}

func TestDumpTable(t *testing.T) {
	if cjdns == nil {
		t.Log("Admin interface not connected.")
		return
	}

	table := cjdns.DumpTable(0)
	if len(table) == 0 {
		t.Fatal("Routing table was not dumped properly.")
	}
}
