// +build windows

package transmission

import (
	"log"
	"net"
)

var tcpPort = "tcp"
var address = "localhost:27015"

// Tunnel function on Windows family systems to connect to default peertalk port
func Tunnel() net.Conn {
	conn, err := net.Dial(tcpPort, address)
	if err != nil {
		log.Println("[TUNNEL-ERROR] : Unable to connect to port.")
	}
	return conn
}
