package cjdngo

import (
	"encoding/json"
	"io/ioutil"
)

// Conf is a struct for cjdroute.conf files. It exports all values.
type Conf struct {
	PrivateKey                  string          `json:"privateKey"`                  //the private key for this node (keep it safe)
	PublicKey                   string          `json:"publicKey"`                   //the public key for this node
	IPv6                        string          `json:"ipv6"`                        //this node's IPv6 address as (derived from publicKey)
	AuthorizedPasswords         []AuthPass      `json:"authorizedPasswords"`         //authorized passwords
	Admin                       AdminBlock      `json:"admin"`                       //information for RCP server
	Interfaces                  InterfacesBlock `json:"interfaces"`                  //interfaces for the switch core
	Router                      RouterBlock     `json:"router"`                      //configuration for the router
	ResetAfterInactivitySeconds int             `json:"resetAfterInactivitySeconds"` //remove cryptoauth sessions after this number of seconds
	PidFile                     string          `json:"pidFile,omitempty"`           //the file to write the PID to, if enabled (disabled by default)
	//BUG(DuoNoxSol): Need to add 'security' block
	Version int `json:"version"` //the internal config file version (mostly unused)
}

//AuthPass is a struct containing a authorization password for connecting peers.
type AuthPass struct {
	Name     string `json:"name,omitempty"`     //the username or real name of the authenticated peer
	Location string `json:"location,omitempty"` //the geographical location of the authenticated peer
	IPv6     string `json:"ipv6,omitempty"`     //the IPv6 used by this peer
	Password string `json:"password"`           //the password for incoming authorization
}

type AdminBlock struct {
	Bind     string `json:"bind"`     //the port to bind the RCP server to
	Password string `json:"password"` //the password for the RCP server
}

type InterfacesBlock struct {
	UDPInterface UDPInterfaceBlock `json:"UDPInterface"`
}

type UDPInterfaceBlock struct {
	Bind      string       `json:"bind"`      //the port to bind the UDP interface to
	ConnectTo ConnectBlock `json:"connectTo"` //the list of peers to connect to
}

type ConnectBlock struct {
	//map here?
}

type RouterBlock struct {
	Interface InterfaceBlock `json:"interface"` //interface used for connecting to the cjdns network
}

type InterfaceBlock struct {
	Type      string `json:"type"`                //the type of interface
	TunDevice string `json:"tunDevice,omitempty"` //the persistent interface to use for cjdns (not usually used)
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
