package cjdngo

import (
	"encoding/json"
	"io/ioutil"
)

// Conf is a struct for cjdroute.conf files. It exports all values.
type Conf struct {
	PrivateKey                  string          //the private key for this node (keep it safe)
	PublicKey                   string          //the public key for this node
	IPv6                        string          //this node's IPv6 address as (derived from publicKey)
	AuthorizedPasswords         []AuthPass      //authorized passwords
	Admin                       AdminBlock      //information for RCP server
	Interfaces                  InterfacesBlock //interfaces for the switch core
	Router                      RouterBlock     //configuration for the router
	ResetAfterInactivitySeconds int             //remove cryptoauth sessions after this number of seconds
	PidFile                     string          //the file to write the PID to, if enabled (disabled by default)
	Version                     int             //the internal config file version (mostly unused)
}

//AuthPass is a struct containing a authorization password for connecting peers.
type AuthPass struct {
	Name string //the username or real name of the authenticated peer
	Location string //the geographical location of the authenticated peer
	IPv6 string //the IPv6 used by this peer
	Password string //the password for incoming authorization
}

type AdminBlock struct {
	Bind     string //the port to bind the RCP server to
	Password string //the password for the RCP server
}

type InterfacesBlock struct {
	UDPInterface UDPInterfaceBlock
}

type UDPInterfaceBlock struct {
	Bind string //the port to bind the UDP interface to
	//connectTo []connectBlock //the list of peers to connect to
}

type RouterBlock struct {
	Interface InterfaceBlock //interface used for connecting to the cjdns network
}

type InterfaceBlock struct {
	Type      string //the type of interface
	TunDevice string //the persistent interface to use for cjdns (not usually used)
}

func ReadConf(path string) (*Conf, error) {
	var conf Conf

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(file, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}


func WriteConf(path string, conf Conf) error {
	b, err := json.MarshalIndent(conf, "", "\t")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0666)
}
