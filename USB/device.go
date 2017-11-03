package USB

import (

)

type (
    
    USBDeviceInfo struct{
        DeviceID uint16
        DeviceUUID string
    }
    
    // Delegate methods for USBDevice, if any ios Device is connected or disconnected
    USBDeviceDelegate interface {
        DeviceDidConnect(USBDeviceInfo)
        DeviceDidDisconnect(USBDeviceInfo)
    }
    
    // Enum for the status of an USBDevice, if its attached to the PC or not
    USBDeviceStatus uint8
)

const(
    ATTACHED USBDeviceStatus = 1 + iota
    DETACHED
)

var connectedDevices map[uint16]USBDeviceInfo

func (newDevice USBDeviceInfo) addToList() {
    connectedDevices[newDevice.DeviceID] = newDevice
}

func (exitingDevice USBDeviceInfo) removeFromList() {
    delete(connectedDevices, exitingDevice.DeviceID)
}
