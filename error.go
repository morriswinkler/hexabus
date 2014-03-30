package hexabus


import "fmt"


type Error struct {
	id int       // error id
	msg string   // error message
	err error   // additional error message
}

func (e Error ) Error() string {
	if e.err != nil {
		return fmt.Sprintf("Error:%d %s %s ", e.id, e.msg, e.err)
	} else {
		return fmt.Sprintf("Error:%d %s", e.id, e.msg)
	}
}

const (
	// Internal errors with id and error message

	// encoder/decoder errors
	ERR_BINWRITE_ID, ERR_BINWRITE_MSG = 20, "binary.Write failed:"
	ERR_BINREAD_ID, ERR_BINREAD_MSG = 21, "binary.READ failed:"
	ERR_MAXSTRBUFF_ID, ERR_MAXSTRBUFF_MSG = 22, "strings can't exeed 127 bytes:" 
	ERR_BYTESIZE_ID, ERR_BYTESIZE_MSG = 23, "16, 65 bytes length are allowed:"
	ERR_PAYLOAD_ID, ERR_PAYLOAD_MSG = 24, "unsuported payload type:"
	ERR_BOOLTYPE_ID, ERR_BOOLTYPE_MSG = 25, "bool can only be 0x00 or 0x01:"
	ERR_HXBDTYPE_ID, ERR_HXBDTYPE_MSG = 26, "unknown hexabus data type:"
	ERR_CRCFAILED_ID, ERR_CRCFAILED_MSG = 27, "checksum mismatch:"

)
