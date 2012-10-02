package cjdngo

import (
	"encoding/json"
	"io/ioutil"
	"regexp"
)

// Conf is a struct for cjdroute.conf files. It exports all values.
type Conf struct {
	Name                        string          `json:"name,omitempty"`              //the username or real name of the person running this node
	Location                    string          `json:"location,omitempty"`          //the geographical location of this node
	PrivateKey                  string          `json:"privateKey"`                  //the private key for this node (keep it safe)
	PublicKey                   string          `json:"publicKey"`                   //the public key for this node
	IPv6                        string          `json:"ipv6"`                        //this node's IPv6 address as (derived from publicKey)
	AuthorizedPasswords         []AuthPass      `json:"authorizedPasswords"`         //authorized passwords
	Admin                       AdminBlock      `json:"admin"`                       //information for RCP server
	Interfaces                  InterfacesBlock `json:"interfaces"`                  //interfaces for the switch core
	Router                      RouterBlock     `json:"router"`                      //configuration for the router
	ResetAfterInactivitySeconds int             `json:"resetAfterInactivitySeconds"` //remove cryptoauth sessions after this number of seconds
	PidFile                     string          `json:"pidFile,omitempty"`           //the file to write the PID to, if enabled (disabled by default)
	Security                    interface{}     `json:"security"`                    //block to contain that strange security formatting
	Version                     int             `json:"version"`                     //the internal config file version (mostly unused)
	//BUG(DuoNoxSol): the Security block is not fully supported
}

//AuthPass is a struct containing a authorization password for connecting peers.
type AuthPass struct {
	Name     string `json:"name,omitempty"`     //the username or real name of the authenticated peer node's owner
	Location string `json:"location,omitempty"` //the geographical location of the authenticated peer node
	IPv6     string `json:"ipv6,omitempty"`     //the IPv6 used by this peer node
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
	Name      string `json:"name,omitempty"`     //the username or real name of the peer node's owner
	Location  string `json:"location,omitempty"` //the geographical location of the peer node
	IPv6      string `json:"ipv6,omitempty"`     //the IPv6 address of the peer node
	Password  string `json:"password"`           //the password to connect to the peer node
	PublicKey string `json:"publicKey"`          //the peer node's public key
}

type RouterBlock struct {
	Interface InterfaceBlock `json:"interface"` //interface used for connecting to the cjdns network
}

type InterfaceBlock struct {
	Type      string `json:"type"`                //the type of interface
	TunDevice string `json:"tunDevice,omitempty"` //the persistent interface to use for cjdns (not usually used)
}

func ReadConf(path string) (*Conf, error) {
	var conf Conf //create a variable to store the conf file in

	file, err := ioutil.ReadFile(path) //read the file from the disk
	if err != nil {
		return nil, err //return the error if necessary
	}

	b := stripComments(file) //strip out all comments, multi and single line

	err = json.Unmarshal(b, &conf) //parse the JSON into a Conf object
	if err != nil {
		return nil, err //return the error if necessary
	}
	return &conf, nil //and return the Conf object
}

func WriteConf(path string, conf Conf) error {
	b, err := json.MarshalIndent(conf, "", "    ") //create a []byte with the JSON, indented by four spaces
	if err != nil {
		return err //return the error if necessary
	}

	return ioutil.WriteFile(path, b, 0666) //and return error from that operation
}

//stripComments replaces all C-style comments (prefixed with "//" and inside "/* */") with empty strings. This is necessary in parsing JSON files that contain them.
//Returns b without comments.
func stripComments(b []byte) []byte {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		panic(err)
	}
	//multiComment, err := regexp.Compile("/\*.*\*/")
	//if err != nil {
	//	panic(err)
	//}

	out := regComment.ReplaceAllLiteral(b, nil)

	return out
}
