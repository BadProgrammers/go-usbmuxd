package USB

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
	"github.com/SoumeshBanerjee/go-usbmuxd/transmission"
	"howett.net/plist"
	"net"
)

func Listen(delegate USBDeviceDelegate) net.Conn {
	// start a tunnel here, and then send the listen frame to that connected socket
	conn, err := transmission.Tunnel()
	if err != nil {
		panic("[USB-ERROR-CONN-1] : Can't establish connection with the USB")
	}
	go frameParser(conn, delegate)

	// send a listen request to usbmuxd daemon socket
	conn.Write(sendListenRequestToSocket())

	return conn
}

func frameParser(conn net.Conn, delegate USBDeviceDelegate) {
	chunk := make([]byte, 2500000)
	for {
		n, err := conn.Read(chunk)
		if err != nil {
			panic("[USB-ERROR-READ-1] : Unable to read data stream from the USB channel")
		}
		// initial check for message type
		var data frames.USBGenericACKFrame
		decoder := plist.NewDecoder(bytes.NewReader(chunk[16:n]))
		decoder.Decode(&data)

		if data.MessageType == "Result" && data.Number != 0 {
			panic("Some error encountered while communication via USB, try again!")
		}
		if data.MessageType != "Result" {
			var data frames.USBDeviceAttachedDetachedFrame
			decoder = plist.NewDecoder(bytes.NewReader(chunk[16:n]))
			decoder.Decode(&data)
			if data.MessageType == "Attached" {
				delegate.DeviceDidConnect(data)
			} else if data.MessageType == "Detached" {
				delegate.DeviceDidDisconnect(data)
			} else {
				delegate.DidReceiveError(errors.New("[USB-ERROR-LISTEN-2] : Unable to parse the response"), string(chunk[16:n]))
			}
		}
	}
}

func sendListenRequestToSocket() []byte {
	// constructing body
	body := &frames.USBListenRequestFrame{
		MessageType:         "Listen",
		ProgName:            "go-usbmuxd",
		ClientVersionString: "1.0.0",
	}
	bodyBuffer := &bytes.Buffer{}
	encoder := plist.NewEncoder(bodyBuffer)
	encoder.Encode(body)

	//constructing header
	headerBuffer := make([]byte, 16)
	binary.LittleEndian.PutUint32(headerBuffer, uint32(bodyBuffer.Len()+16))
	headerBuffer[4] = byte(1)
	headerBuffer[8] = byte(8)
	headerBuffer[12] = byte(1)

	requestBuffer := append(headerBuffer, bodyBuffer.Bytes()...)
	return requestBuffer
}
