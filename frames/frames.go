package frames

type (

	// When we want to listen for any new USB device or device removed
	USBListenRequestFrame struct {
		MessageType         string `plist:"MessageType"`
		ClientVersionString string `plist:"ClientVersionString"`
		ProgName            string `plist:"ProgName"`
	}

	// Its a frame model for generic response after we send listen or connect
	// Number == 0 {OK}, Number == 1 {Device not connected anymore}, Number == 2 {Port not available}, Number == 5 {IDK}
	USBGenericACKFrame struct {
		MessageType string `plist:"MessageType"`
		Number      int    `plist:"Number"`
	}

	// Model for USB connect or disconnect frame
	USBDeviceAttachedDetachedFrame struct {
		MessageType string                               `plist:"MessageType"`
		DeviceID    int                                  `plist:"DeviceID"`
		Properties  USBDeviceAttachedPropertiesDictFrame `plist:"Properties"`
	}

	//Model for USB attach properties
	USBDeviceAttachedPropertiesDictFrame struct {
		ConnectionSpeed int    `plist:"ConnectionSpeed"`
		ConnectionType  string `plist:"ConnectionType"`
		DeviceID        int    `plist:"DeviceID"`
		LocationID      int    `plist:"LocationID"`
		ProductID       int    `plit:"ProductID"`
		SerialNumber    string `plit:"SerialNumber"`
	}

	// Model for connect frame to a specific port in a connected device
	USBConnectRequestFrame struct {
		MessageType         string `plist:"MessageType"`
		ClientVersionString string `plist:"ClientVersionString"`
		ProgName            string `plist:"ProgName"`
		DeviceID            int    `plist:"DeviceID"`
		PortNumber          int    `plist:"PortNumber"`
	}
)
