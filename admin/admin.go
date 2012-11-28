//Python tools for accessing the CJDNS admin interface are available [here](https://github.com/cjdelisle/cjdns/tree/master/contrib/python#cjdnspy).
//cjdngo/admin is an attempt to port those utilities to Go. Other functions and utilities will be made available.

package admin

import (
	"github.com/zeebo/bencode"
	"net"
	//"time"
)

//This is the type which wraps a CJDNS admin interface. It is initialized by Connect().
type CJDNS struct {
	address  string //Address to connect to
	password string //Admin interface password, used to connect
	port     string //Port the admin interface is bound to
	cookie   string //Cookie as given by the interface on connection
}

//The CJDNS admin interface has a large number of functions. These are the string constants used to invoke them.
const (
	CommandPing = "ping" //Check if the admin server is running
	CommandAuth = "auth" //Use authentication
)

//The CJDNS admin interface sometimes responds with particular strings to indicate statuses, such as "pong" to indicate that it is running.
const (
	StatusPingOK = "pong" //Response from CommandPing
)

func Connect(address, password, port string) (cjdns *CJDNS, err error) {
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

//Wraps the command and arguments in a map[string]interface{}, then uses the given Conn to encode them directly to the wire.
func send(conn net.Conn, command, cookie, hash string, args map[string]string) {
	//Exit if the command is not given.
	if command == "" {
		return
	}
	//Otherwise, create the map which will be used
	//to encode the message.
	message := make(map[string]interface{})
	if cookie != "" && hash != "" {
		//If there is authentication involved,
		//then use "aq". Otherwise, "q".
		message["q"] = CommandAuth
		message["aq"] = command
		message["cookie"] = cookie
		message["hash"] = hash
	} else {
		message["q"] = command
	}

	if args != nil {
		message["args"] = args
	}
	bencode.NewEncoder(conn).Encode(message)
}

func receive(conn net.Conn) (response map[string]interface{}) {
	bencode.NewDecoder(conn).Decode(&response)
	return
}

//Sends the admin server a ping, and returns true if it gets the expected response.
func (cjdns *CJDNS) Ping() (status bool) {
	conn, err := net.Dial("tcp", cjdns.address+":"+cjdns.port)
	if err != nil {
		return
	}
	defer conn.Close()

	send(conn, CommandPing, "", "", nil)

	response := receive(conn)

	if response["q"] != StatusPingOK {
		return
	}
	return true
}
