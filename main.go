package main

import (
	"github.com/SoumeshBanerjee/go-usbmuxd/USB"
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
	"github.com/SoumeshBanerjee/go-usbmuxd/transmission"
	"io"
	"log"
	"os"
)

// some global vars
var connectHandle USB.ConnectedDevices
var port = 29173
var pluggedUSBDevices map[int]frames.USBDeviceAttachedDetachedFrame
var connectedUSB int // only stores the device id
var scanningInstance USB.Scan
var self USBDeviceDelegate

func main() {
	// inti section
	connectedUSB = -1
	pluggedUSBDevices = map[int]frames.USBDeviceAttachedDetachedFrame{}
	scanningInstance = USB.Scan{}
	self = USBDeviceDelegate{}

	// logger
	logFile, err := os.OpenFile("kusb_ios.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	defer logFile.Close()
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	// create a USB.Listen(USBDeviceDelegate) instance. Pass a delegate to resolve the attached and detached callbacks
	// then on device added save ot to array/ map and send connect to a port with proper tag
	listenConnection := USB.Listen(transmission.Tunnel(), self)
	defer listenConnection.Close()

	// connect to a random usb device, if Number == 0 then
	connectHandle = USB.ConnectedDevices{Delegate: self, Connection: transmission.Tunnel()}
	defer connectHandle.Connection.Close()

	// scan defer
	defer scanningInstance.Stop()

	// run loop
	select {}
}

// MARK: - USB Delegate Methods
type USBDeviceDelegate struct{}

func (usb USBDeviceDelegate) USBDeviceDidPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// usb has been plugged DO: startScanning
	log.Printf("[USB-INFO] : Device Plugged %s ID: %d\n", frame.Properties.SerialNumber, frame.DeviceID)
	pluggedUSBDevices[frame.DeviceID] = frame
	scanningInstance.Start(&connectHandle, frame, port)
}
func (usb USBDeviceDelegate) USBDeviceDidUnPlug(frame frames.USBDeviceAttachedDetachedFrame) {
	// usb has been unplugged
	// stop scan
	log.Printf("[USB-INFO] : Device UnPlugged %s ID: %d\n", pluggedUSBDevices[frame.DeviceID].Properties.SerialNumber, frame.DeviceID)
	delete(pluggedUSBDevices, frame.DeviceID)
	scanningInstance.Stop()
}
func (usb USBDeviceDelegate) USBDidReceiveErrorWhilePluggingOrUnplugging(err error, stringResponse string) {
	// plug or unplug error
	// stop scan
	if stringResponse != "" {
		//some unresolved message came
		//TODO - Implement some resolver to understand message received
	}
	log.Println("[USB-EM-1] : Some error encountered wile pluging and unpluging. ", err.Error())
	scanningInstance.Stop()
}
func (usb USBDeviceDelegate) USBDeviceDidSuccessfullyConnect(device USB.ConnectedDevices, deviceID int, toPort int) {
	// successfully connected to the port mentioned
	// stop the scan
	connectedUSB = deviceID
	scanningInstance.Stop()
}
func (usb USBDeviceDelegate) USBDeviceDidFailToConnect(device USB.ConnectedDevices, deviceID int, toPort int, err error) {
	// error while communication in the socket
	// start scan
	connectedUSB = -1
	pluggedDeviceID := getFirstPluggedDeviceId()
	if pluggedDeviceID != -1 {
		scanningInstance.Start(&connectHandle, pluggedUSBDevices[pluggedDeviceID], port)
	}

}
func (usb USBDeviceDelegate) USBDeviceDidReceiveData(device USB.ConnectedDevices, deviceID int, messageTAG uint32, data []byte) {
	// received data from the device
	log.Println(string(data))
	//device.SendData(data[20:], 106)
}
func (usb USBDeviceDelegate) USBDeviceDidDisconnect(devices USB.ConnectedDevices, deviceID int, toPort int) {
	// socket disconnect
	// start scan
	connectedUSB = -1
	pluggedDeviceID := getFirstPluggedDeviceId()
	if pluggedDeviceID != -1 {
		scanningInstance.Start(&connectHandle, pluggedUSBDevices[pluggedDeviceID], port)
	}
}

// MARK - helper functions here
func getFirstPluggedDeviceId() int {
	var deviceID int = -1
	for deviceID, _ = range pluggedUSBDevices {
		break
	}
	return deviceID
}
