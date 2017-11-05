// +build windows
package transmission

import (
	"net"
)

func Tunnel() net.Conn {
	conn, err := net.Dial("tcp", "localhost:27015")
	if err != nil {
		log.Println("[TUNNEL-ERROR] : Unable to connect to port.")
	}
	return conn
}
