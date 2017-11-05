package main

import (
	"fmt"
	"github.com/SoumeshBanerjee/go-usbmuxd/USB"
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
    "encoding/binary"
)

var connectHandle USB.ConnectedDevices
var port = 2345
func main() {
 
	// create a USB.Listen(USBDeviceDelegate) instance. Pass a delegate to resolve the attached and detached callbacks
	// then on device added save ot to array/ map and send connect to a port with proper tag
	self := USBDeviceDelegate{}
	listenConnection := USB.Listen(self)
	defer listenConnection.Close()

	// connect to a random usb device, if Number == 0 then
	connectHandle = USB.ConnectedDevices{Delegate: self}
	
	// run loop
	select {}
}

// MARK: - USB Delegate Methods
type USBDeviceDelegate struct{}

func (usb USBDeviceDelegate) USBDeviceDidPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// callback will arrive here for new plug in usb device
	fmt.Println("DIDConnect : " + frame.MessageType)
	connectHandle.Connect(frame, port)

}
func (usb USBDeviceDelegate) USBDeviceDidUnPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// disconnect call will come here
	fmt.Println("didDISConnect : " + frame.MessageType)
}
func (usb USBDeviceDelegate) USBDidReceiveErrorWhilePluggingOrUnplugging(err error, stringResponse string) {
	if err != nil {
		panic(err)
	}
}
func (usb USBDeviceDelegate) USBDeviceDidSuccessfullyConnect(device USB.ConnectedDevices, deviceID int, toPort int) {
    fmt.Println("Device Connected Successfully")
    for i:=0; i<=10; i++ {
        if device.Connection != nil {
            data := helper("hi")
            _, err := device.Connection.Write(data)
            if err != nil {
                fmt.Println(err.Error())
            }
        }
    }
}
func (usb USBDeviceDelegate) USBDeviceDidFailToConnect(device USB.ConnectedDevices, deviceID int, toPort int, err error) {
    fmt.Println("Failed to connect" + err.Error())
}
func (usb USBDeviceDelegate) USBDeviceDidReceiveData(device USB.ConnectedDevices, deviceID int, data []byte) {
    fmt.Println(data)
}
func (usb USBDeviceDelegate) USBDeviceDidDisconnect(devices USB.ConnectedDevices, deviceID int, toPort int)  {
    fmt.Println("Socket Disconnect")
}

func helper(string string) []byte {
    bodyBuffer := []byte(string)
    
    headerBuffer := make([]byte, 20)
    
    binary.BigEndian.PutUint32(headerBuffer[:4], 1)
    binary.BigEndian.PutUint32(headerBuffer[4:8], 101)
    binary.BigEndian.PutUint32(headerBuffer[8:12], 0)
    binary.BigEndian.PutUint32(headerBuffer[12:16], uint32(len(bodyBuffer)+4))
    binary.BigEndian.PutUint32(headerBuffer[16:], uint32(len(bodyBuffer)))
    
    return append(headerBuffer, bodyBuffer...)
}
