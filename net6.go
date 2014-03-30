package hexabus


import (
	"net"
	"time"
	"regexp"
	"fmt"
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
        
	n, err := conn.Read(readbuf)
	
        if err != nil {
		if opErr, ok := err.(net.Error); ok && !opErr.Timeout() {
		return  err
		}
        }
	
	if n > 0 {
		err = checkHeader(readbuf[:n])
		if err != nil {
			return err
		}
		ptype, err := PacketType(readbuf[:n])
		if err != nil {
			return err
		}
		if ptype == PTYPE_ERROR {
			ep := ErrorPacket{} 
			ep.Decode(readbuf[:n])
			return Error{id:ERR_ERRPACKET_ID, msg:ERR_ERRPACKET_MSG + fmt.Sprintf("%d",ep.Error)}
		}
	}
	return nil
}

