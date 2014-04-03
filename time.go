package hexabus

import (
	"bytes"
	"encoding/binary"
)

// struct to hold DTYPE_TIMESTAMP
type Timestamp struct {
	TotalSeconds uint32
}

// decodes the payload from a timestamp packet into a Timestamp structure
// takes as argument a Packet.Data interface{}
func (t *Timestamp) Decode(data interface{}) (err error) {
	buf := bytes.NewBuffer(data.([]byte))
	err = binary.Read(buf, binary.BigEndian, t)
	if err != nil {
		return err
	}
	return nil
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
func (d *DateTime) Decode(data interface{}) (err error) {
	buf := bytes.NewBuffer(data.([]byte))
	err = binary.Read(buf, binary.BigEndian, d)
	if err != nil {
		return err
	}
	return nil
}
