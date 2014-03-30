package hexabus


import (
	"net"
	"time"
	"io/ioutil"
)


func send(packet []byte, address string) ([]byte, error) {
	
	
	conn, err := net.DialTimeout("udp6", address, time.Duration(NET_TIMEOUT)*time.Second)
	if err != nil {
		return nil, err
	}
	
	_, err = conn.Write(packet)
	if err != nil {
		return nil, err
	}
	
	result, err := ioutil.ReadAll(conn)
	if err != nil {
		return nil, err
	}
	
	return result, nil
}


