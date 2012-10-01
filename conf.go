package cjdngo

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Conf struct {
	privateKey string //the private key for this node (keep it safe)
	publicKey string //the public key for this node
	ipv6 string //this node's IPv6 address as (derived from publicKey)
	authorizedPasswords []authPass //authorized passwords
	admin adminBlock //information for RCP server
	interfaces interfacesBlock //interfaces for the switch core
	router routerBlock //configuration for the router
	resetAfterInactivitySeconds int //remove cryptoauth sessions after this number of seconds
	pidFile string //the file to write the PID to, if enabled (disabled by default)
}
	
type authPass struct {
	password string //the password for incoming authorization
	//add "name" and "location" fields?
}

type adminBlock struct {
	bind string //the port to bind the RCP server to
	password string //the password for the RCP server
}

type interfacesBlock struct {
	UDPInterface UDPInterfaceBlock
}

type UDPInterfaceBlock struct {
	bind string //the port to bind the UDP interface to
	//connectTo []connectBlock //the list of peers to connect to
}

type routerBlock struct {
	interfac interfaceBlock //interface used for connecting to the cjdns network
}

type interfaceBlock struct {
	typ string //the type of interface
	tunDevice string //the persistent interface to use for cjdns (not usually used)
}

func ReadConf(path string) *Conf {
	var conf Conf
	
	file, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	
	err = json.Unmarshal(file, &conf)
	if err != nil {
		panic(err)
	}
	return &conf
}

func WriteConf(path string, conf Conf) error {
	 b, err := json.Marshal(conf)
	 if err != nil {
	 	return err
	 }
	 
	 err = ioutil.WriteFile(path, b)
}
