// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris

package transmission

import (
	"net"
)

var osFamily = "unix"
var tcpFile = "/var/run/usbmuxd"

// Tunnel USB function for Unix family
func Tunnel() (net.Conn, error) {
	conn, err := net.Dial(osFamily, tcpFile)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
