package main

import (
	"fmt"
	"github.com/SoumeshBanerjee/go-usbmuxd/USB"
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
)

func main() {

	// create a USB.Listen(USBDeviceDelegate) instance. Pass a delegate to resolve the attached and detached callbacks
	// then on device added save ot to array/ map and send connect to a port with proper tag
	listenConnection := USB.Listen(USBDeviceDelegate{})
	defer listenConnection.Close()

	// connect to a random usb device, if Number == 0 then

	// run loop
	select {}
}

func byteSwap(val int) int {
	return ((val & 0xFF) << 8) | ((val >> 8) & 0xFF)
}

// MARK: - USB Delegate Methods
type USBDeviceDelegate struct{}

func (usb USBDeviceDelegate) DeviceDidConnect(frame frames.USBDeviceAttachedDetachedFrame) {
	// callback will arrive here for new plug in usb device
	fmt.Println("DIDConnect : " + frame.MessageType)

}
func (usb USBDeviceDelegate) DeviceDidDisconnect(frame frames.USBDeviceAttachedDetachedFrame) {
	// disconnect call will come here
	fmt.Println("didDISConnect : " + frame.MessageType)
}
func (usb USBDeviceDelegate) DidReceiveError(err error, stringResponse string) {
	if err != nil {
		panic(err)
	}
}
