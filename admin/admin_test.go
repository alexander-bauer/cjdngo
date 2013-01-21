package admin

import (
	"github.com/SashaCrofter/cjdngo"
	"testing"
)

var (
	cjdns *CJDNS
	table []*Route
)

func TestConnect(t *testing.T) {
	conf, err := cjdngo.ReadConf("/etc/cjdroute.conf")
	if err != nil || len(conf.Admin.Password) == 0 {
		// This is not related to the test.
		t.Log(err)
		t.Fatal("Could not read the config file. This is not related to cjdngo/admin.")
	}

	cjd, err := Connect("127.0.0.1", "11234", conf.Admin.Password)
	if err != nil {
		t.Fatal(err)
	}

	pingOK, authOK := cjd.Ping()
	if !authOK {
		t.Fatal("Server authentication failed.")
	} else if !pingOK {
		t.Fatal("Server did not respond to ping.")
	}
	cjdns = cjd
}

func TestDumpTable(t *testing.T) {
	if cjdns == nil {
		t.Log("Admin interface not connected; skipping test.")
		return
	}

	table = cjdns.DumpTable(-1)
	if len(table) == 0 {
		table = nil
		t.Fatal("Routing table was not dumped properly.")
	}
	t.Log("Number of routes is", len(table))
}

func TestPeerFilter(t *testing.T) {
	if cjdns == nil {
		t.Log("Admin interface not connected; skipping test.")
		return
	}
	if table == nil {
		t.Log("Routing table could not be dumped; skipping test.")
		return
	}

	peers := cjdns.Peers(1) // Retrieve direct peers
	if peers == nil {
		t.Fatal("Peers were returned nil.")
	}
	t.Logf("Number of direct peers: %d\n", len(peers))

	extendedPeers := cjdns.Peers(2) // Peers of peers
	if extendedPeers == nil {
		t.Fatal("Extended peers were returned nil.")
	}
	t.Logf("Number of extended peers: %d\n", len(extendedPeers))
}
