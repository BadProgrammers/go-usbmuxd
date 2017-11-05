package USB

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
	"howett.net/plist"
	"io"
	"net"
    "log"
)

type (
	ConnectedDeviceDelegate interface {
		USBDeviceDidSuccessfullyConnect(device ConnectedDevices, deviceID int, toPort int)
		USBDeviceDidFailToConnect(device ConnectedDevices, deviceID int, toPort int, err error)
		USBDeviceDidReceiveData(device ConnectedDevices, deviceID int, messageTAG uint32, data []byte)
		USBDeviceDidDisconnect(devices ConnectedDevices, deviceID int, toPort int)
	}
	ConnectedDevices struct {
		Delegate   ConnectedDeviceDelegate
		Connection net.Conn
	}
)

func (device ConnectedDevices) Connect(conn net.Conn, frame frames.USBDeviceAttachedDetachedFrame, port int) net.Conn {
	device.Connection.Close()
	device.Connection = conn
	go connectFrameParser(device.Connection, frame.DeviceID, port, device)

	// now send a connect request to the device
	device.Connection.Write(sendConnectRequestToSocket(frame.DeviceID, port))

	return device.Connection
}

func (device ConnectedDevices) SendData(data []byte, messageTagType uint32) {
	// create a 20byte standard header that's used by peertalk to parse tag and other info
	headerBuffer := make([]byte, 20)

	// preparing the header
	binary.BigEndian.PutUint32(headerBuffer[:4], 1)
	binary.BigEndian.PutUint32(headerBuffer[4:8], messageTagType)
	binary.BigEndian.PutUint32(headerBuffer[8:12], 0)
	binary.BigEndian.PutUint32(headerBuffer[12:16], uint32(len(data)+4))
	binary.BigEndian.PutUint32(headerBuffer[16:], uint32(len(data)))

	if device.Connection != nil {
		_, err := device.Connection.Write(append(headerBuffer, data...))
		if err != nil {
		    log.Println(err.Error())
		}
	}
}

func byteSwap(val int) int {
	return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}

func sendConnectRequestToSocket(deviceID int, toPort int) []byte {
	// constructing body
	body := &frames.USBConnectRequestFrame{
		DeviceID:            deviceID,
		PortNumber:          byteSwap(toPort),
		MessageType:         "Connect",
		ClientVersionString: "1.0.0",
		ProgName:            "go-usbmuxd",
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

func connectFrameParser(conn net.Conn, deviceID int, toPort int, device ConnectedDevices) {
	chunk := make([]byte, 2500000)
	for {
		n, err := conn.Read(chunk)
		if err != nil {
			if err != io.EOF {
				panic("[USB-ERROR-iCONNECT-1] : Unable to read data stream from the USB channel")
			}
			device.Delegate.USBDeviceDidDisconnect(device, deviceID, toPort)
			break
		}
		// initial check for message type
		var data frames.USBGenericACKFrame
		decoder := plist.NewDecoder(bytes.NewReader(chunk[16:n]))
		decoder.Decode(&data)
		if data.MessageType == "Result" && data.Number == 0 {
			device.Delegate.USBDeviceDidSuccessfullyConnect(device, deviceID, toPort)
		} else if data.MessageType == "Result" && data.Number != 0 {
			var errorMessage string
			switch data.Number {
			case 2:
				// Device Disconnected
				errorMessage = "Unable to connect to Device, might be issue with the cable or turned off"
			case 3:
				// Port isn't available/ busy
				errorMessage = "Port you're requesting is unavailable"
			case 5:
				// UNKNOWN Error
				errorMessage = "[IDK]: Malformed request received in the device"
			}
			device.Delegate.USBDeviceDidFailToConnect(device, deviceID, toPort, errors.New(errorMessage))
		}
		if data.MessageType != "Result" {
			// parse the TAG and other relevant header info
			headerBuffer := chunk[:16]
			//log.Println(binary.BigEndian.Uint32(headerBuffer[0:4]))
			//log.Println(binary.BigEndian.Uint32(headerBuffer[4:8]))
			//log.Println(binary.BigEndian.Uint32(headerBuffer[8:12]))
			//log.Println(binary.BigEndian.Uint32(headerBuffer[12:16]))
			//log.Println(binary.BigEndian.Uint32(headerBuffer[16:20]))
			device.Delegate.USBDeviceDidReceiveData(device, deviceID, binary.BigEndian.Uint32(headerBuffer[4:8]), chunk[:n])
		}
	}
}
