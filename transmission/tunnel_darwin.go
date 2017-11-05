// +build darwin dragonfly freebsd linux nacl netbsd openbsd solaris
package transmission

import (
	"net"
)

func Tunnel() (net.Conn, error) {
	conn, err := net.Dial("unix", "/var/run/usbmuxd")
	if err != nil {
		return nil, err
	}
	return conn, nil
}