// +build windows

package transmission

import (
	"net"
)

var connectionProtocol = "tcp"
var connectionSocket = "127.0.0.1:27015"

// Tunnel USB function for Windows family
func Tunnel() (net.Conn, error) {
	conn, err := net.Dial(connectionProtocol, connectionSocket)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
