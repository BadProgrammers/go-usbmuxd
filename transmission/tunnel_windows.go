// +build windows

package transmission

import (
	"net"
)

func Tunnel() (net.Conn, error) {
	conn, err := net.Dial("tcp", "27015")
	if err != nil {
		return nil, err
	}
	return conn, nil
}


