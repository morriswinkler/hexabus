package hexabus

import "fmt"
import "bytes"
import "encoding/binary"

func addHeader(packet []byte) {
	packet[0], packet[1], packet[2], packet[3] = HEADER0, HEADER1, HEADER2, HEADER3
}

// set datatype and encode payload in bytes
func encData(packet []byte, data interface{}) []byte {
	switch data := data.(type) {
	case bool:
		packet[10] = DTYPE_BOOL
		if data == true {
			packet = append(packet, TRUE)
		} else {
			packet = append(packet, FALSE)
		}
	case uint8:
		packet[10] = DTYPE_UINT8
		packet = append(packet, data)
	case uint32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = DTYPE_UINT32
		packet = append(packet, buf.Bytes()...)
	// DateTime: holds DTYPE_DATETIME data ; needs testing
	case DateTime:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = DTYPE_DATETIME
		packet = append(packet, buf.Bytes()...)
	case float32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = DTYPE_FLOAT
		packet = append(packet, buf.Bytes()...)
	case string:
		// TODO: check if you can send smaller string length then 128 bytes
		// might be the same case as in 16BYTES and 66BYTES
		if len(data) > STRING_PACKET_MAX_BUFFER_LENGTH {
			panic(fmt.Errorf("max string length 127 exeeded for string: %s", data))
		} else {
			// TODO: check if 0 termination in string is right that way
			packet[10] = DTYPE_128STRING
			packet = append(packet, data...)
			packet = append(packet, byte(0))
		}
		// TIMESTAMP: intended for type syscall.Sysinfo_t.Uptime not working
	case Timestamp:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			panic(fmt.Errorf("binary.Write failed:", err))
		}
		packet[10] = DTYPE_TIMESTAMP
		packet = append(packet, buf.Bytes()...)
	case []byte:
		// there are only 16, 66 bytes long byte packets, they where both added to
		// serve a uniq purpos, bytes with variable length is planned in the next protokoll version or so
		if len(data) == 16 {
			packet[10] = DTYPE_16BYTES
			packet = append(packet, data...)
		} else if len(data) == 65 {
			packet[10] = DTYPE_66BYTES
			packet = append(packet, data...)
		} else {
			panic(fmt.Errorf("only 16 or 66 bytes of raw data are allowed length %d is not supported", len(data)))
		}
	default:
		packet[10] = DTYPE_UNDEFINED
		panic(fmt.Errorf("unsupported payload type: %T", data))
	}

	return packet

}

func decData(data []byte, dtype byte) (ret_data interface{}) {
	switch dtype {
	case DTYPE_BOOL:
		if data[0] == 0x01 {
			ret_data = true
		} else if data[0] == 0x00 {
			ret_data = false
		} else {
			panic(fmt.Errorf("data type bool cant be: %T", data))
		}
	case DTYPE_UINT8:
		ret_data = uint8(data[0])
	case DTYPE_UINT32:
		var v uint32
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			panic(fmt.Errorf("binary.Read failed:", err))
		}
		ret_data = v
	case DTYPE_DATETIME:
		var v DateTime
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			panic(fmt.Errorf("binary.Read failed:", err))
		}
		ret_data = v
	case DTYPE_FLOAT:
		var v float32
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			panic(fmt.Errorf("binary.Read failed:", err))
		}
		ret_data = v
	case DTYPE_128STRING:
		ret_data = string(data[0:len(data)-1])
	case DTYPE_TIMESTAMP:
		var v Timestamp
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			panic(fmt.Errorf("binary.Read failed:", err))
		}
		ret_data = v
	case DTYPE_16BYTES:
		ret_data = data
	case DTYPE_66BYTES:
		ret_data = data
	default:
		panic(fmt.Errorf("unknown hexabus data type: %x", dtype))
	}

	return 
}  

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

func checkCRC(packet []byte) ( result bool) {
	crc_c := crc16(packet[:len(packet)-2])
	crc_r := binary.BigEndian.Uint16(packet[len(packet)-2:])
	if crc_c == crc_r {
		return true
	} else {
		panic(fmt.Errorf("checksum %d and %d do not match", crc_c, crc_r))  
		return false
	}
}

// struct to hold DTYPE_TIMESTAMP
type Timestamp struct {
	TotalSeconds uint32
}

// decodes the payload from a timestamp packet into a Timestamp structure
// takes as argument a Packet.Data interface{}
func (t *Timestamp) Decode(data interface{}) {
	buf := bytes.NewBuffer(data.([]byte))
	err := binary.Read(buf, binary.BigEndian, t)
	if err != nil {
		panic(fmt.Errorf("binary.Write failed:", err))
	}
}

// struct to hold DTYPE_DATETIME
type DateTime struct {
	Hours     uint8
	Minutes   uint8
	Seconds   uint8
	Day       uint8
	Month     uint8
	Year      uint16
	DayOfWeek uint8
}

// decodes the payload from a datetime packet into a DateTime structure
// takes as argument a Packet.Data interface{}
func (d *DateTime) Decode(data interface{}) {
	buf := bytes.NewBuffer(data.([]byte))
	err := binary.Read(buf, binary.BigEndian, d)
	if err != nil {
		panic(fmt.Errorf("binary.Write failed:", err))
	}
}

type ErrorPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte // 1 byte flags
	Error byte // 1 byte error code
}

func (p *ErrorPacket) Encode() []byte {
	packet := make([]byte, 7)
	addHeader(packet)
	packet[4] = PTYPE_ERROR
	packet[5] = p.Flags
	packet[6] = p.Error
	packet = addCRC(packet)
	return packet
}

func (p *ErrorPacket) Decode(packet []byte) {
	if checkCRC(packet) {
		p.Flags = packet[5]
		p.Error = packet[6]
	}
}

type InfoPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte        // 1 byteflags
	Eid   uint32      // 4 bytes endpoint id
	Dtype byte        // 1 byte data type
	Data  interface{} // ... bytes payload, size depending on datatype
}

func (p *InfoPacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_INFO
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet[10] = p.Dtype
	packet = encData(packet, p.Data)
	packet = addCRC(packet)
	return packet
}

func (p *InfoPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = decData(packet[11 : len(packet)-2], packet[10])
}

type QueryPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte   // flags
	Eid   uint32 // endpoint id
}

func (p *QueryPacket) Encode() []byte {
	packet := make([]byte, 12)
	addHeader(packet)
	packet[4] = PTYPE_QUERY
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addCRC(packet)
	return packet
}

func (p *QueryPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
}

type WritePacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte        // flags
	Eid   uint32      // endpoint id
	Dtype byte        // data type
	Data  interface{} // payload, size depending on datatype
}

func (p *WritePacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_WRITE
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = encData(packet, p.Data)
	packet = addCRC(packet)
	return packet
}

func (p *WritePacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = packet[11 : len(packet)-2]

}
