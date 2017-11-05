package USB

import (
    "github.com/SoumeshBanerjee/go-usbmuxd/frames"
    "bytes"
    "howett.net/plist"
    "encoding/binary"
    "github.com/SoumeshBanerjee/go-usbmuxd/transmission"
    "net"
    "errors"
    "io"
)

type (
    ConnectedDeviceDelegate interface {
        USBDeviceDidSuccessfullyConnect(device ConnectedDevices, deviceID int, toPort int)
        USBDeviceDidFailToConnect(device ConnectedDevices, deviceID int, toPort int, err error)
        USBDeviceDidReceiveData(device ConnectedDevices, deviceID int, data []byte)
        USBDeviceDidDisconnect(devices ConnectedDevices, deviceID int, toPort int)
    }
    ConnectedDevices struct {
        Delegate ConnectedDeviceDelegate
        Connection net.Conn
    }
)

func (device ConnectedDevices) Connect (frame frames.USBDeviceAttachedDetachedFrame, port int) net.Conn {
    conn, err := transmission.Tunnel()
    if err != nil {
        panic("[USB-CONNECT-ERROR-1] : Unable to connect to USB")
    }
    device.Connection = conn
    go connectFrameParser(conn, frame.DeviceID, port, device)
    
    // now send a connect request to the device
    conn.Write(sendConnectRequestToSocket(frame.DeviceID, port))
    
    return conn
}

func byteSwap(val int) int {
    return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}

func sendConnectRequestToSocket(deviceID int, toPort int) []byte {
    // constructing body
    body := &frames.USBConnectRequestFrame{
        DeviceID: deviceID,
        PortNumber: byteSwap(toPort),
        MessageType: "Connect",
        ClientVersionString: "1.0.0",
        ProgName: "go-usbmuxd",
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
            if err != io.EOF{
                panic("[USB-ERROR-iCONNECT-1] : Unable to read data stream from the USB channel")
            }
            device.Delegate.USBDeviceDidDisconnect(device, deviceID, toPort)
            break
        }
        // initial check for message type
        var data frames.USBGenericACKFrame
        decoder := plist.NewDecoder(bytes.NewReader(chunk[16:n]))
        decoder.Decode(&data)
        if data.MessageType == "Result" && data.Number == 0{
            device.Delegate.USBDeviceDidSuccessfullyConnect(device, deviceID, toPort)
        }else if data.MessageType == "Result" && data.Number != 0 {
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
            //headerBuffer := chunk[:16]
            //fmt.Println(binary.BigEndian.Uint32(headerBuffer[0:4]))
            //fmt.Println(binary.BigEndian.Uint32(headerBuffer[4:8]))
            //fmt.Println(binary.BigEndian.Uint32(headerBuffer[8:12]))
            //fmt.Println(binary.BigEndian.Uint32(headerBuffer[12:16]))
            //fmt.Println(binary.BigEndian.Uint32(headerBuffer[16:20]))
            device.Delegate.USBDeviceDidReceiveData(device, deviceID, chunk[:n])
        }
    }
}