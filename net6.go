package hexabus

import (
	"fmt"
	"net"
	"regexp"
	"time"
)

// Defaults used by the network communication.
const (
	// hexabus default port
	PORT = "61616"

	// package transmit timeout
	NET_TIMEOUT = 3
)

// structure to hold all EID's of a hexabus device and its capabilities
type EID struct {
	Eid      uint32 // Eid
	Dtype    byte   // data type
	Desc     string // description
	Writable bool   // writeable
}

// Switch IPv6 address, quantity Example:
// eids, err := hexabus.QueryEids("[fafa::50:c4ff:fe04:8390]", 32)
func QueryEids(address string, eid_qty uint16) ([]EID, error) {
	eid_mask := []uint16{}
	eid_descriptors := []uint16{}
	eid_map := []EID{}

	// find all EID's in eid_qty that are 0 or can be multiplyed by 32
	for i := uint16(0); i < eid_qty; i = i + 32 {
		eid_descriptors = append(eid_descriptors, uint16(i))
	}

	// build eid_mask to check what EID's are available
	for _, descriptor := range eid_descriptors {
		pq := QueryPacket{FLAG_NONE, uint32(descriptor)}
		result, err := pq.Send(address)
		if err != nil {
			return nil, err
		}
		// change byte order to read LSB first
		result = result[len(result)-6 : len(result)-2]
		result = []byte{result[3], result[2], result[1], result[0]}

		for _, bit := range result {
			eid_mask = append(eid_mask, uint16((bit)&1))
			eid_mask = append(eid_mask, uint16((bit>>1)&1))
			eid_mask = append(eid_mask, uint16((bit>>2)&1))
			eid_mask = append(eid_mask, uint16((bit>>3)&1))
			eid_mask = append(eid_mask, uint16((bit>>4)&1))
			eid_mask = append(eid_mask, uint16((bit>>5)&1))
			eid_mask = append(eid_mask, uint16((bit>>6)&1))
			eid_mask = append(eid_mask, uint16((bit>>7)&1))
		}
	}

	// query all availabel EID's and build a map of struc EID
	for eid, avlb := range eid_mask {
		if avlb == 1 {
			result, err := EpQueryPacket{0, uint32(eid)}.Send(address)
			if err != nil {
				return nil, err
			}
			pei := EpInfoPacket{}
			err = pei.Decode(result)
			if err != nil {
				return nil, err
			}
			var data interface{}
			if pei.Dtype != DTYPE_UINT8 {
				data = uint8(1)
			} else {
				data = uint32(1)
			}

			// check if endpoint is writable
			err = WritePacket{FLAG_NONE, uint32(eid), DTYPE_UNDEFINED, data}.Send(address)
			if err == Error(0x02) {
				eid_map = append(eid_map, EID{uint32(eid), pei.Dtype, pei.Data.(string), false})
			} else if err == Error(0x04) {
				eid_map = append(eid_map, EID{uint32(eid), pei.Dtype, pei.Data.(string), true})
			} else {
				return nil, err
			}
		}
	}

	return eid_map, nil
}

func (p QueryPacket) Send(address string) ([]byte, error) {

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

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(NET_TIMEOUT * time.Second)))
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

func (p WritePacket) Send(address string) error {

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

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(NET_TIMEOUT * time.Second)))
	if err != nil {
		return err
	}

	// Register time before writting the packet
	timeBeforeWrite := time.Now()

	// Write the packet
	var wn int
	wn, err = conn.Write(packet)
	if err != nil {
		return err
	}

	// Check the wn variable
	fmt.Printf("wn: %v\n", wn)

	// Register the time after writing and report
	timeAfterWrite := time.Now()
	fmt.Printf("Write took %v seconds\n",
		timeAfterWrite.Sub(timeBeforeWrite).Seconds())

	// Register time before reading the packet
	timeBeforeRead := time.Now()

	// Read the packet
	n, err := conn.Read(readbuf)

	fmt.Printf("n: %v\n", n)

	// Register time after read and report
	timeAfterRead := time.Now()
	fmt.Printf("Reading took %v seconds\n",
		timeAfterRead.Sub(timeBeforeRead).Seconds())

	// Check errors
	if err != nil {
		if opErr, ok := err.(net.Error); ok && !opErr.Timeout() {
			return err
		}
	}

	if n > 0 {
		err = checkCRC(readbuf[:n])
		if err != nil {
			return err
		}
		err = checkHeader(readbuf[:n])
		if err != nil {
			return err
		}
		ptype, err := PacketType(readbuf[:n])
		if err != nil {
			return err
		}

		fmt.Printf("Packet type %v\n", ptype)

		if ptype == PTYPE_ERROR {
			ep := ErrorPacket{}
			ep.Decode(readbuf[:n])
			return Error(ep.Error)
		}
	}
	return nil
}

func (p EpQueryPacket) Send(address string) ([]byte, error) {

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

	err = conn.SetReadDeadline(time.Now().Add(time.Duration(NET_TIMEOUT * time.Second)))
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
