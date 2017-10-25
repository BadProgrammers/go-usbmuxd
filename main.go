package main

import (
	"fmt"
	"encoding/binary"
	"bytes"
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
    <string>node-usbmux</string>
    <key>ProgName</key>
    <string>node-usbmux</string>
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
	fmt.Println(req_buf)

}


//func reader(r io.Reader) {
//	buf := make([]byte, 1024)
//	for {
//		print("here")
//		n, err := r.Read(buf[:])
//		if err != nil {
//			return
//		}
//		println("Client got:", string(buf[0:n]))
//	}
//}