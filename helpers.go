package hexabus

import (
	"bytes"
	"encoding/binary"
)

// add packet header
func addHeader(packet []byte) {
	packet[0], packet[1], packet[2], packet[3] = HEADER0, HEADER1, HEADER2, HEADER3
}

// check if a received packet header is valid
func checkHeader(packet []byte) error {
	if packet[0] == HEADER0 && packet[1] == HEADER1 && packet[2] == HEADER2 && packet[3] == HEADER3 {
		return nil
	}
	return Error{id: ERR_WRONGHEADER_ID, msg: ERR_WRONGHEADER_MSG}
}

// set datatype and encode payload in bytes
func encData(packet []byte, data interface{}) ([]byte, error) {
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
			return nil, Error{id: ERR_BINWRITE_ID, msg: ERR_BINWRITE_MSG, err: err}
		}
		packet[10] = DTYPE_UINT32
		packet = append(packet, buf.Bytes()...)
	// DateTime: holds DTYPE_DATETIME data ; needs testing
	case DateTime:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return nil, Error{id: ERR_BINWRITE_ID, msg: ERR_BINWRITE_MSG, err: err}
		}
		packet[10] = DTYPE_DATETIME
		packet = append(packet, buf.Bytes()...)
	case float32:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return nil, Error{id: ERR_BINWRITE_ID, msg: ERR_BINWRITE_MSG, err: err}
		}
		packet[10] = DTYPE_FLOAT
		packet = append(packet, buf.Bytes()...)
	case string:
		// TODO: check if you can send smaller string length then 128 bytes
		// might be the same case as in 16BYTES and 66BYTES
		if len(data) > STRING_PACKET_MAX_BUFFER_LENGTH {
			return nil, Error{id: ERR_MAXSTRBUFF_ID, msg: ERR_MAXSTRBUFF_MSG + data}
		} else {
			// TODO: check if 0 termination in string is right that way
			packet[10] = DTYPE_128STRING
			packet = append(packet, data...)
			packet = append(packet, byte(0))

			for len(packet[11:]) < 128 {
				packet = append(packet, byte(0))
			}
		}
		// TIMESTAMP: intended for type syscall.Sysinfo_t.Uptime not working
	case Timestamp:
		buf := new(bytes.Buffer)
		err := binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return nil, Error{id: ERR_BINWRITE_ID, msg: ERR_BINWRITE_MSG, err: err}
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
			return nil, Error{id: ERR_BYTESIZE_ID, msg: ERR_BYTESIZE_MSG + string(len(data))}
		}
	default:
		packet[10] = DTYPE_UNDEFINED
		return nil, Error{id: ERR_PAYLOAD_ID, msg: ERR_PAYLOAD_MSG}
	}

	return packet, nil

}

// decode received Data payload
func decData(data []byte, dtype byte) (interface{}, error) {
	var ret_data interface{}
	switch dtype {
	case DTYPE_BOOL:
		if data[0] == 0x01 {
			ret_data = true
		} else if data[0] == 0x00 {
			ret_data = false
		} else {
			return nil, Error{id: ERR_BOOLTYPE_ID, msg: ERR_BOOLTYPE_MSG}
		}
	case DTYPE_UINT8:
		ret_data = uint8(data[0])
	case DTYPE_UINT32:
		var v uint32
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			return nil, Error{id: ERR_BINREAD_ID, msg: ERR_BINREAD_MSG, err: err}
		}
		ret_data = v
	case DTYPE_DATETIME:
		var v DateTime
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			return nil, Error{id: ERR_BINREAD_ID, msg: ERR_BINREAD_MSG, err: err}
		}
		ret_data = v
	case DTYPE_FLOAT:
		var v float32
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			return nil, Error{id: ERR_BINREAD_ID, msg: ERR_BINREAD_MSG, err: err}
		}
		ret_data = v
	case DTYPE_128STRING:
		if len(data) != 128 {
			return nil, Error{id: ERR_MAXSTRBUFF_ID, msg: ERR_MAXSTRBUFF_MSG + string(data)}
		}
		end := bytes.IndexByte(data, 0x00)
		if end != -1 {
			ret_data = string(data[0:end])
		} else {
			return nil, Error{id: ERR_STRNOTERM_ID, msg: ERR_STRNOTERM_MSG}
		}
	case DTYPE_TIMESTAMP:
		var v Timestamp
		buf := bytes.NewBuffer(data)
		err := binary.Read(buf, binary.BigEndian, &v)
		if err != nil {
			return nil, Error{id: ERR_BINREAD_ID, msg: ERR_BINREAD_MSG, err: err}
		}
		ret_data = v
	case DTYPE_16BYTES:
		if len(data) != 16 {
			return nil, Error{id: ERR_BYTESIZE_ID, msg: ERR_BYTESIZE_MSG}
		}
		ret_data = data
	case DTYPE_66BYTES:
		if len(data) != 65 {
			return nil, Error{id: ERR_BYTESIZE_ID, msg: ERR_BYTESIZE_MSG}
		}
		ret_data = data
	default:
		return nil, Error{id: ERR_HXBDTYPE_ID, msg: ERR_HXBDTYPE_MSG}
	}

	return ret_data, nil
}

func PacketType(packet []byte) (ptype byte, err error) {
	switch packet[4] {
	case PTYPE_ERROR:
		ptype = PTYPE_ERROR
	case PTYPE_INFO:
		ptype = PTYPE_INFO
	case PTYPE_QUERY:
		ptype = PTYPE_QUERY
	case PTYPE_WRITE:
		ptype = PTYPE_WRITE
	case PTYPE_EPINFO:
		ptype = PTYPE_EPINFO
	case PTYPE_EPQUERY:
		ptype = PTYPE_EPQUERY
	default:
		return 0xff, Error{id: ERR_UNKNOWNPTYPE_ID, msg: ERR_UNKNOWNPTYPE_MSG}
	}

	return ptype, nil
}
