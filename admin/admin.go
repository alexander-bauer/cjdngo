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

//The CJDNS admin interface has a large number of functions. These are the string constants used to invoke them.
const (
	CommandAuth   = "auth"   //Use authentication
	CommandPing   = "ping"   //Check if the admin server is running
	CommandCookie = "cookie" //Request a cookie from the server
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

//Wraps the command and arguments in a map[string]interface{}, then uses the given Conn to encode them directly to the wire.
func (cjdns *CJDNS) Send(conn net.Conn, command string, args map[string]string) (response map[string]interface{}) {
	//Exit if the command is not given.
	if command == "" {
		return
	}
	//Otherwise, create the map which will be used
	//to encode the message.
	message := make(map[string]interface{})

	if args != nil {
		message["args"] = args
	}

	if cjdns.cookie != "" && cjdns.password != "" {
		//If there is authentication involved,
		//then use "aq". Otherwise, "q".
		message["q"] = CommandAuth
		message["aq"] = command
		message["cookie"] = cjdns.cookie
		message["hash"] = cjdns.password + cjdns.cookie //as specified

		//Prepare the hash
		m, err := bencode.EncodeString(message)
		if err != nil {
			return
		}

		hash := sha256.New()
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
