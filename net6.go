package hexabus


import (
	"net"
	"time"
	"regexp"
)


func (p *QueryPacket) Send(address string) ([]byte, error) {
	
	packet := p.Encode()
	
	// check if port is set otherwhise append default hexabus port
	var validPort = regexp.MustCompile(`:[0-9]{1,5}$`)
	if !validPort.MatchString(address) {
		address += ":" + PORT
	}
	readbuf := make([]byte, 152)
        conn, err := net.DialTimeout("udp6", address, time.Duration(NET_TIMEOUT)*time.Second)
        if err != nil {
                return nil, err
        }

        _, err = conn.Write(packet)
        if err != nil {
                return nil, err
        }

        n, err := conn.Read(readbuf)
        if err != nil {
                return nil, err
        }

        return readbuf[:n], nil
}

func (p *WritePacket) Send(address string) (error) {
	
	packet, err := p.Encode()
	if err != nil {
		return err
	}

	// check if port is set otherwhise append default hexabus port
	var validPort = regexp.MustCompile(`:[0-9]{1,5}$`)
	if !validPort.MatchString(address) {
		address += ":" + PORT
	}
	readbuf := make([]byte, 152)
        conn, err := net.DialTimeout("udp6", address, time.Duration(NET_TIMEOUT)*time.Second)
        if err != nil {
                return err
        }

        _, err = conn.Write(packet)
        if err != nil {
                return err
        }

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(NET_TIMEOUT * time.Second)))
	if err != nil {
		return  err
	}
        
	_, err = conn.Read(readbuf)
        if err != nil {
                return  err
        }
	
	return nil
	
}

