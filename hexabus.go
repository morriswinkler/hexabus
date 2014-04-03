// Hexabus is the go implementation of the Hexabus packet format used by Hexabus devices.
// For more info see https://github.com/mysmartgrid/hexabus.
package hexabus

// Hexabus Error Packet
// if a packet os malformed or doesn't the EID is not properly used it will return
// a error of type ERR_UNKNOWNEID, ERR_WRITEREADONLY, ERR_CRCFAILED, ERR_DATATYPE
// or ERR_INVALID_VALUE
type ErrorPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte // 1 byte flags
	Error byte // 1 byte error code
}

// encoder for Error Packet
func (p *ErrorPacket) Encode() []byte {
	packet := make([]byte, 7)
	addHeader(packet)
	packet[4] = PTYPE_ERROR
	packet[5] = p.Flags
	packet[6] = p.Error
	packet = addCRC(packet)
	return packet
}

// decoder for Error Packet
func (p *ErrorPacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Error = packet[6]
	return nil
}

// Hexabus Info Packet
// Info Packets are send after a Query Packet or Broadcasted every n seconds
type InfoPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte        // 1 byteflags
	Eid   uint32      // 4 bytes endpoint id
	Dtype byte        // 1 byte data type
	Data  interface{} // ... bytes payload, size depending on datatype
}

// encoder for Info Packet
func (p *InfoPacket) Encode() (packet []byte, err error) {
	packet = make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_INFO
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet[10] = p.Dtype
	packet, err = encData(packet, p.Data)
	if err != nil {
		return nil, err
	}
	packet = addCRC(packet)
	return packet, nil
}

// decoder for Info Packet
func (p *InfoPacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data, err = decData(packet[11:len(packet)-2], packet[10])
	if err != nil {
		return err
	}
	return nil
}

// Hexabus Query Packet
// used to query endpoints to return an Info Packet
type QueryPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte   // flags
	Eid   uint32 // endpoint id
}

// encoder for Query Packet
func (p *QueryPacket) Encode() []byte {
	packet := make([]byte, 10)
	addHeader(packet)
	packet[4] = PTYPE_QUERY
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addCRC(packet)
	return packet
}

// decoder for Query Packet
func (p *QueryPacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	return nil
}

// Hexabus Write Packet
// used to set a writeable entpoint id to a certain value, there is no response for that o
// other than an Error Packet on fail
type WritePacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte        // flags
	Eid   uint32      // endpoint id
	Dtype byte        // data type
	Data  interface{} // payload, size depending on datatype
}

// encoder for Write Packet
func (p *WritePacket) Encode() (packet []byte, err error) {
	packet = make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_WRITE
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet, err = encData(packet, p.Data)
	if err != nil {
		return nil, err
	}
	packet = addCRC(packet)
	return packet, nil
}

// decoder for Write Packet
func (p *WritePacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	p.Dtype = packet[10]
	p.Data, err = decData(packet[11:len(packet)-2], packet[10])
	if err != nil {
		return err
	}
	return nil
}

// Hexabus Endpoint Info Package
// this is a response packet to Endpoint Query and holds the EID Dtype and a
// message that describes that endpoint in Data
type EpInfoPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte        // 1 byteflags
	Eid   uint32      // 4 bytes endpoint id
	Dtype byte        // 1 byte data type
	Data  interface{} // ... bytes payload, size depending on datatype
}

// encoder for Endpoint Info Packet
func (p *EpInfoPacket) Encode() (packet []byte, err error) {
	packet = make([]byte, 11)
	addHeader(packet)
	packet[4] = PTYPE_EPINFO
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet[10] = p.Dtype
	packet, err = encData(packet, p.Data)
	if err != nil {
		return nil, err
	}
	packet = addCRC(packet)
	return packet, nil
}

// decoder for Endpoint Info Packet
func (p *EpInfoPacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	// Endpoin Info Packets have the datatype of the endpoint that was queried
	// with Endpoint Query, so for example a Relay will anther with type Bool
	// to turn it on and off
	p.Dtype = packet[10]
	// we need to set the data type here to DTYPE_128STRING since this is allwayes
	// a description of the queried endpoint
	// don't use the datatype from p.Dtype
	p.Data, err = decData(packet[11:len(packet)-2], DTYPE_128STRING)
	if err != nil {
		return err
	}
	return nil
}

// Hexabus Endpoint Query Packet
// used to query endpoints for a data ytpe and a short description
type EpQueryPacket struct {
	// 4 bytes header
	// 1 byte packet type
	Flags byte   // flags
	Eid   uint32 // endpoint id
}

// encoder for Endpoint Query Packet
func (p *EpQueryPacket) Encode() []byte {
	packet := make([]byte, 10)
	addHeader(packet)
	packet[4] = PTYPE_EPQUERY
	packet[5] = p.Flags
	packet[6], packet[7], packet[8], packet[9] = uint8(p.Eid>>24), uint8(p.Eid>>16), uint8(p.Eid>>8), uint8(p.Eid&0xff)
	packet = addCRC(packet)
	return packet
}

// decoder for Endpoint Query Packet
func (p *EpQueryPacket) Decode(packet []byte) (err error) {
	err = checkHeader(packet)
	if err != nil {
		return err
	}
	err = checkCRC(packet)
	if err != nil {
		return err
	}
	p.Flags = packet[5]
	p.Eid = uint32(uint8(packet[6])>>24 + uint8(packet[7])>>16 + uint8(packet[8])>>8 + uint8(packet[9])&0xff)
	return nil
}
