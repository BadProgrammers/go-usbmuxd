package main

import (
	"fmt"
	"encoding/binary"
	"github.com/SoumeshBanerjee/go-usbmuxd/transmission"
	"io"
	"time"
)

type header struct {
	len uint32
	version uint32
	request uint32
	tag uint32
}

func main() {
	str := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>MessageType</key>
    <string>Listen</string>
    <key>ClientVersionString</key>
    <string>1.0.0</string>
    <key>ProgName</key>
    <string>go-usbmuxd</string>
  </dict>
</plist>`


	buf := []byte(str)

	_ = header{len: uint32(len(buf)+16), version:1, request:8, tag:1}

	header_buffer := make([]byte, 16)




	binary.LittleEndian.PutUint32(header_buffer, uint32(len(buf)+16))
	header_buffer[4] = byte(1)
	header_buffer[8] = byte(8)
	header_buffer[12] = byte(1)



	req_buf := append(header_buffer, buf...)


	conn, err := transmission.Tunnel()
	if err!=nil {
		fmt.Println(err)
	}

	go reader(conn)

	_, err = conn.Write(req_buf)

	if err!=nil {
		fmt.Println("Writing Error: ", err)
	}

	// run loop

	for {
		time.Sleep(1)
	}

}


func reader(r io.Reader) {
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			return
		}
		fmt.Println(string(buf[16:n]))
	}
}