// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris
package transmission

import (
	"log"
	"net"
)

func Tunnel() net.Conn {
	conn, err := net.Dial("unix", "/var/run/usbmuxd")
	if err != nil {
		log.Println("[TUNNEL-ERROR] : Unable to connect to port.")
	}
	return conn
}
