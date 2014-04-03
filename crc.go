package hexabus

import (
	"encoding/binary"
)

// calculate crc16 variant
// this code was translated from a php snippet found on http://www.lammertbies.nl/forum/viewtopic.php?t=1253
func crc16(packet []byte) uint16 {
	var crc uint16
	for _, v := range packet {
		crc = crc ^ uint16(v)
		for y := 0; y < 8; y++ {
			if (crc & 0x001) == 0x0001 {
				crc = (crc >> 1) ^ 0x8408
			} else {
				crc = crc >> 1
			}
		}
	}
	// in the original Kermit implementation the two crc bytes are swaped
	// i commented that since this is not kermit
	// hexabus uses the same CRC as contiki os , type ???
	//lb := (crc & 0xff00) >> 8
	//hb := (crc & 0x00ff) << 8
	//crc = hb | lb
	return crc
}

// add checksum
func addCRC(packet []byte) []byte {
	crc := crc16(packet)
	packet = append(packet, uint8(crc>>8), uint8(crc&0xff))
	return packet
}

// validate checksum from a received hexabus package
func checkCRC(packet []byte) (err error) {
	crc_c := crc16(packet[:len(packet)-2])
	crc_r := binary.BigEndian.Uint16(packet[len(packet)-2:])
	if crc_c == crc_r {
		return nil
	} else {
		return Error(0xa5)
	}
}
