package USB

import (
	"../frames"
)

type (
	// DeviceDelegate defines Delegate methods for USBDevice, if any ios Device is plugged or unplugged
	DeviceDelegate interface {
		USBDeviceDidPlug(frames.USBDeviceAttachedDetachedFrame)
		USBDeviceDidUnPlug(frames.USBDeviceAttachedDetachedFrame)
		USBDidReceiveErrorWhilePluggingOrUnplugging(error, string)
	}
)
