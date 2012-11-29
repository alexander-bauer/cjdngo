//Python tools for accessing the CJDNS admin interface are available [here](https://github.com/cjdelisle/cjdns/tree/master/contrib/python#cjdnspy).
//cjdngo/admin is an attempt to port those utilities to Go. Other functions and utilities will be made available.

package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/zeebo/bencode"
	"io"
	"net"
	//"time"
)

//This is the type which wraps a CJDNS admin interface. It is initialized by Connect().
type CJDNS struct {
	address  string //Address to connect to
	password string //Admin interface 64-bit hash, used to connect
	port     string //Port the admin interface is bound to
	cookie   string //Cookie as given by the interface on connection
}

//A Route is a single entry in the CJDNS routing table. Each entry has preciesly one IPv6 address, one path, one unitless link quality number, and one version number.
type Route struct {
	IP      string `bencode:"ip"`      //Node's IPv6 address
	Path    string `bencode:"path"`    //Routing path to the node
	Link    int64  `bencode:"link"`    //The link quality (unitless)
	Version int64  `bencode:"version"` //The node's version (primarily unused)
}

//The CJDNS admin interface has a large number of functions. These are the string constants used to invoke them.
const (
	CommandAuth      = "auth"                //Use authentication
	CommandPing      = "ping"                //Check if the admin server is running
	CommandCookie    = "cookie"              //Request a cookie from the server
	CommandDumpTable = "NodeStore_dumpTable" //Dump the routing table
)

//The CJDNS admin interface sometimes responds with particular strings to indicate statuses, such as "pong" to indicate that it is running.
const (
	StatusPingOK = "pong" //Response from CommandPing
)

var (
	NoPingError         = errors.New("admin interface not responding")
	NoCookieError       = errors.New("admin interface did not offer cookie")
	AuthenticationError = errors.New("admin interface rejected password")
)

func Connect(address, port, password string) (cjdns *CJDNS, err error) {
	if len(address) == 0 || len(port) == 0 || len(password) == 0 {
		return nil, errors.New("not enough arguments to Connect()")
	}
	cjdns = &CJDNS{
		address:  address,
		password: password,
		port:     port,
	}
	up := cjdns.Ping()
	if !up {
		//If the interface doesn't respond,
		//return the NoPingError.
		return nil, NoPingError
	}

	//Get a cookie from the interface.
	cjdns.cookie = cjdns.Cookie()
	if len(cjdns.cookie) == 0 {
		//If the server did not offer a
		//cookie, return the NoCookieError.
		return nil, NoCookieError
	}
	return
}

//Wrapper for CJDNS.Ping, which does not require a password. It is useful to check if the CJDNS admin server is running, without necessarily having access to a configuration file. It will return true if the server is up.
func Ping(address, port string) (status bool) {
	cjdns := &CJDNS{
		address: address,
		port:    port,
	}
	return cjdns.Ping()
}

//Wrapper for CJDNS.Cookie, which does not require a password. It returns the string that the admin server generates.
func Cookie(address, port string) (cookie string) {
	cjdns := &CJDNS{
		address: address,
		port:    port,
	}
	return cjdns.Cookie()
}

//Wraps the command and arguments in a map[string]interface{}, then uses the given Conn to encode them directly to the wire. It sends authorization if it is supplied in the given CJDNS.
func (cjdns *CJDNS) Send(conn net.Conn, command string, args map[string]interface{}) (response map[string]interface{}) {
	//Exit if the command is not given.
	if command == "" {
		return
	}
	//Otherwise, create the map which will be used
	//to encode the message.
	message := make(map[string]interface{})

	if cjdns.cookie != "" && cjdns.password != "" {
		//If there is authentication involved,
		//then use "aq". Otherwise, "q".

		hash := sha256.New()
		hash.Write([]byte(cjdns.password + cjdns.cookie))

		message["aq"] = command
		message["args"] = args
		message["cookie"] = cjdns.cookie
		message["hash"] = hex.EncodeToString(hash.Sum(nil)) //as specified
		message["q"] = CommandAuth

		//Prepare the hash
		m, err := bencode.EncodeString(message)
		if err != nil {
			return
		}

		hash = sha256.New()
		hash.Write([]byte(m))
		message["hash"] = hex.EncodeToString(hash.Sum(nil))
	} else {
		message["q"] = command
	}

	m, err := bencode.EncodeString(message)
	if err == nil {
		io.WriteString(conn, m)
	}

	bencode.NewDecoder(conn).Decode(&response)
	return
}

//Wraps the net.Dial() function to reach an admin server. It is the caller's responsibility to close the connection.
func (cjdns *CJDNS) Dial() (conn net.Conn, err error) {
	return net.Dial("tcp", cjdns.address+":"+cjdns.port)
}

//Sends the admin server a ping, and returns true if it gets the expected response.
func (cjdns *CJDNS) Ping() (status bool) {
	conn, err := cjdns.Dial()
	if err != nil {
		return
	}
	defer conn.Close()

	response := cjdns.Send(conn, CommandPing, nil)

	if response["q"] != StatusPingOK {
		return
	}
	return true
}

//Asks the admin server for a cookie, and returns the resultant string.
func (cjdns *CJDNS) Cookie() (cookie string) {
	conn, err := cjdns.Dial()
	if err != nil {
		return
	}
	defer conn.Close()

	response := cjdns.Send(conn, CommandCookie, nil)

	return response["cookie"].(string)
}

//Retrieves the desired page of the routing table from the admin server. Generally, one wants page 0. Requires authorization.
func (cjdns *CJDNS) DumpTable(page int) (table []*Route) {
	conn, err := cjdns.Dial()
	if err != nil {
		return
	}
	defer conn.Close()

	args := make(map[string]interface{}, 1)
	args["page"] = page

	response := cjdns.Send(conn, CommandDumpTable, args)
	rawTable := response["routingTable"].([]interface{})
	table = make([]*Route, len(rawTable))
	for i := range rawTable {
		item := rawTable[i].(map[string]interface{})
		table[i] = &Route{
			IP:      item["ip"].(string),
			Path:    item["path"].(string),
			Link:    item["link"].(int64),
			Version: item["version"].(int64),
		}
	}
	return
}
