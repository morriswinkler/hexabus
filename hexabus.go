package hexabus

import "fmt"
import "bytes"
import "encoding/binary"

func addHeader(packet []byte) {
	packet[0], packet[1], packet[2], packet[3] = HEADER0, HEADER1, HEADER2, HEADER3
}

func addData(packet []byte, data interface{}) []byte { 	
	// set datatype
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
		// TODO: check if padding is the intended behavior
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

// calculate crc16 kermit variant
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
	// hexabus uses the same CRC as contiki os , type ???
        //lb := (crc & 0xff00) >> 8
        //hb := (crc & 0x00ff) << 8
        //crc = hb | lb
	return crc
} 

func addCRC(packet []byte) []byte {
	crc := crc16(packet)
	packet = append(packet,uint8(crc>>8), uint8(crc&0xff))
	return packet
}

// struct to hold DTYPE_TIMESTAMP
type Timestamp struct {
	TotalSeconds uint32
}

func (t *Timestamp) Decode(data interface{}) {
        buf := bytes.NewBuffer(data.([]byte))
        err := binary.Read(buf, binary.BigEndian, t)
        if err != nil {
                panic(fmt.Errorf("binary.Write failed:", err))
        }
}

// struct to hold DTYPE_DATETIME data
type DateTime struct {
        Hours uint8
        Minutes uint8
        Seconds uint8
        Day uint8
        Month uint8
        Year uint16
        DayOfWeek uint8
}

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
	Flags byte	 // 1 byte flags 
	Error byte       // 1 byte error code
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
	p.Flags = packet[5]
	p.Error = packet[6]
}

type InfoPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte	  // 1 byteflags 
	Eid uint32  // 4 bytes endpoint id
	Dtype byte	  // 1 byte data type
	Data interface{}   // ... bytes payload, size depending on datatype
}

func (p *InfoPacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_INFO
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet[10] = p.Dtype
	packet = addData(packet, p.Data)
	packet = addCRC(packet)     
	return packet 
}

func (p *InfoPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = packet[11:len(packet)-2]
}

// remove that !!!
func EncodeInfoPacket( flags byte, eid uint32, data interface{}) (p []byte) {
	packet := make([]byte, 11, 141)                                                
	addHeader(packet)
        packet[4] = PTYPE_INFO
        packet[5] = flags
        packet[6], packet[7], packet[8], packet[9] = uint8(eid>>24), uint8(eid>>16), uint8(eid>>8), uint8(eid&0xff)
	fmt.Printf("EID bits: %b, %b, %b, %b\n", packet[6], packet[7], packet[8], packet[9])
	packet = addData(packet, data)                                                
        packet = addCRC(packet)                                                         
        return packet
}
	
type QueryPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte	 // flags 
	Eid uint32       // endpoint id
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

func (p  *QueryPacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
}

type WritePacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte       // flags 
	Eid uint32       // endpoint id
	Dtype byte	 // data type
	Data interface{} // payload, size depending on datatype
}     

func (p *WritePacket) Encode() []byte {
	packet := make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_WRITE
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addData(packet, p.Data)                                                
	packet = addCRC(packet)     
	return packet 
}

func (p *WritePacket) Decode(packet []byte) {
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data = packet[11:len(packet)-2]

}
