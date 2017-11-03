package USB

import (
	"github.com/SoumeshBanerjee/go-usbmuxd/frames"
)

type (
	// Delegate methods for USBDevice, if any ios Device is connected or disconnected
	USBDeviceDelegate interface {
		DeviceDidConnect(frames.USBDeviceAttachedDetachedFrame)
		DeviceDidDisconnect(frames.USBDeviceAttachedDetachedFrame)
		DidReceiveError(error, string)
	}
)

var connectedDevices map[uint16]frames.USBDeviceAttachedDetachedFrame
