// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package transmission

import (
	"log"
	"net"
)

var tcpPort = "unix"
var address = "/var/run/usbmuxd"

// Tunnel function on Unix family systems to connect to default peertalk port
func Tunnel() net.Conn {
	conn, err := net.Dial(tcpPort, address)
	if err != nil {
		log.Println("[TUNNEL-ERROR] : Unable to connect to port.")
	}
	return conn
}
