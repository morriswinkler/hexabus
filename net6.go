package hexabus


import (
	"net"
	"time"
	"regexp"
)


func Send(packet []byte, address string) ([]byte, error) {
	
	// check if port is set otherwhise append default hexabus port
	var validPort = regexp.MustCompile(`:[0-9]{1,5}$`)
	if !validPort.MatchString(address) {
		address += ":" + string(PORT)
	}
	
	readbuf := make([]byte, 10, 100)
        conn, err := net.DialTimeout("udp6", address, time.Duration(NET_TIMEOUT)*time.Second)
        if err != nil {
                return nil, err
        }

        _, err = conn.Write(packet)
        if err != nil {
                return nil, err
        }

        _, err = conn.Read(readbuf)
        if err != nil {
                return nil, err
        }

        return readbuf, nil
}

