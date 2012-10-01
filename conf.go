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
	Version                     int             `json:"version"`                     //the internal config file version (mostly unused)
	//BUG(DuoNoxSol): Need to add 'security' block
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
	Bind      string                `json:"bind"`      //the port to bind the UDP interface to
	ConnectTo map[string]Connection `json:"connectTo"` //maps connection information to peer details, where the Key is the peer's IPv4 address and port (or other connection detail) and the Element contains all of the information about the peer, such as password and public key
}

//Connection describes authentication details for connection to a peer who is serving this node. It is stored in the config file as dependent to a string, such as an IPv4 address (and port,) which is necessary to connect to the peer.
type Connection struct {
	Name      string `json:"name,omitempty"`     //the username or real name of the peer
	Location  string `json:"location,omitempty"` //the geographical location of the peer
	IPv6      string `json:"ipv6,omitempty"`     //the IPv6 address of the peer
	Password  string `json:"password"`           //the password to connect to the peer
	PublicKey string `json:"publicKey"`          //the peer's public key
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
	b, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, b, 0666)
}
