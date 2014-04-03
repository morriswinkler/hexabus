package hexabus


import "fmt"

// Error type structure. 
type Error struct {
	id int       // error id
	msg string   // error message
	err error   // additional error message
}

// Error returns an error type.
// All Errors are passed with an ID and optional an err type from
// other packages.
func (e Error ) Error() string {
	if e.err != nil {
		return fmt.Sprintf("Error:%d %s %s ", e.id, e.msg, e.err)
	} else {
		return fmt.Sprintf("Error:%d %s", e.id, e.msg)
	}
}

// Defaults used by the network communication. 
const (
	// hexabus default port
	PORT = "61616"

	// package transmit timeout
	NET_TIMEOUT = 3
)


// Internal error message.
const (
	// encoder/decoder errors
	ERR_BINWRITE_ID, ERR_BINWRITE_MSG = 20, "binary.Write failed:"
	ERR_BINREAD_ID, ERR_BINREAD_MSG = 21, "binary.READ failed:"
	ERR_MAXSTRBUFF_ID, ERR_MAXSTRBUFF_MSG = 22, "strings can't exeed 127 bytes:" 
	ERR_STRNOTERM_ID, ERR_STRNOTERM_MSG = 23, "string is not 0 terminated:"
	ERR_BYTESIZE_ID, ERR_BYTESIZE_MSG = 24, "16, 65 bytes length are allowed:"
	ERR_PAYLOAD_ID, ERR_PAYLOAD_MSG = 25, "unsuported payload type:"
	ERR_BOOLTYPE_ID, ERR_BOOLTYPE_MSG = 26, "bool can only be 0x00 or 0x01:"
	ERR_HXBDTYPE_ID, ERR_HXBDTYPE_MSG = 27, "unknown hexabus data type:"
	ERR_CRCFAILED_ID, ERR_CRCFAILED_MSG = 28, "checksum mismatch:"

	// network errors
	ERR_WRONGHEADER_ID, ERR_WRONGHEADER_MSG = 40, "wrong packet header:"
	ERR_UNKNOWNPTYPE_ID, ERR_UNKNOWNPTYPE_MSG = 41, "unknown packet type:"
	ERR_ERRPACKET_ID, ERR_ERRPACKET_MSG = 42, "received error packet with value:"
)
