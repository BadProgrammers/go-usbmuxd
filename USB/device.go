package USB

import (
	"../frames"
)

type (
	// Delegate methods for USBDevice, if any ios Device is plugged or unplugged
	USBDeviceDelegate interface {
		USBDeviceDidPlug(frames.USBDeviceAttachedDetachedFrame)
		USBDeviceDidUnPlug(frames.USBDeviceAttachedDetachedFrame)
		USBDidReceiveErrorWhilePluggingOrUnplugging(error, string)
	}
)
