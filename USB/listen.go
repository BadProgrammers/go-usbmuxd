package USB

import (
    "net"
)

func Listen(delegate USBDeviceDelegate) net.Conn {
    // some logic here
    // start a tunnel here, and then send the listen frame to that connected socket
    // if device added
    
    
    
    return nil
}

func frameParser(conn net.Conn)  {
    buff := make([]byte, 5000000000)
    for {
        _, err := conn.Read(buff)
        if err!=nil {
            panic("[USB-ERROR-LISTENING-1] : " + err.Error())
        }
        
        
        // call plist package and parse the message to switch the message type
        
        
    }
}