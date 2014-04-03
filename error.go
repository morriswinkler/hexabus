package hexabus

import "strconv"

// Error type structure to hold err.id. err.msg and optional err.err.
type Error int

// Error returns an error type.
// All Errors are passed with ID and MSG and optional an err type from
// other packages.
func (e Error) Error() string {
	str := error_message[e]
	if str == "" {
		// Please report a bug if this happens
		return "Hexabus error #" + strconv.Itoa(int(e))
	}

	return str
}

var error_message = map[Error]string{
	// Hexanus packet errors
	HXB_ERR_SUCCESS:       "hexabus packet error success",
	HXB_ERR_UNKNOWNEID:    "hexabus packet endpoint does not exist",
	HXB_ERR_WRITEREADONLY: "hexabus packet write on read only eid",
	HXB_ERR_CRCFAILED:     "hexabus packet crc failed",
	HXB_ERR_DATATYPE:      "hexabus packet data type doesn't fit endpoint",
	HXB_ERR_INVALID_VALUE: "hexabus packet value can not be interpreted",

	// internal errors
	ERR_STRBUFF:   "strings must be 127 bytes",
	ERR_STRNOTERM: "string is not 0 terminated",
	ERR_BYTESIZE:  "only bytes with 16 or 65 bit length are allowed",
	ERR_HXBDTYPE:  "unsuported hexbus data type",
	ERR_BOOLTYPE:  "bool can only be 0x00 or 0x01",
	ERR_CRCFAILED: "checksum mismatch",

	// internal network errors
	ERR_WRONGHEADER:  "wrong packet header",
	ERR_UNKNOWNPTYPE: "unknown packet type",
	ERR_ERRPACKET:    "received error packet with value",
}

// Internal error codes.
const (
	// encoder/decoder errors
	ERR_STRBUFF   = 0xa0
	ERR_STRNOTERM = 0xa1
	ERR_BYTESIZE  = 0xa2
	ERR_HXBDTYPE  = 0xa3
	ERR_BOOLTYPE  = 0xa4
	ERR_CRCFAILED = 0xa5

	// network errors
	ERR_WRONGHEADER  = 0xb0
	ERR_UNKNOWNPTYPE = 0xb1
	ERR_ERRPACKET    = 0xb2
)
