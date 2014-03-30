package hexabus

// Constants
const (
	/* Header */

	// The UDP Data of a Hexabus Packet starts with the Bytes 0x48 0x58 0x30 0x43
	// (HX0C) to identify it as a Hexabus Packet

	HEADER0, HEADER1, HEADER2, HEADER3 = 0x48, 0x58, 0x30, 0x43

	/* Boolean values */

	// boolean false
	FALSE = 0x00

	// boolean true
	TRUE = 0x01

	/* Packet types */

	// Hexabus Error Packet
	// An error occured -- check the error code field for more information
	PTYPE_ERROR = 0x00

	// Hexabus Info Packet
	// Endpoint provides information
	PTYPE_INFO = 0x01

	// Hexabus Query Packet
	// Endpoint is requested to provide information
	PTYPE_QUERY = 0x02

	// Hexabus Write Packet
	// Endpoint is requested to set its value
	PTYPE_WRITE = 0x04

	// Hexabus EpInfo Packet
	// Endpoint metadata
	PTYPE_EPINFO = 0x09

	// Hexabus EpQuery Packet
	// Request endpoint metadata
	PTYPE_EPQUERY = 0x09

	/* Flags */

	// Hexabus Flag "No Flag set"
	FLAG_NONE = 0x00

	/* Data types */

	// Hexabus Data Type "No data at all"
	DTYPE_UNDEFINED = 0x00

	// Data type Bool
	DTYPE_BOOL = 0x01

	// Data type uint8
	DTYPE_UINT8 = 0x02

	// Data type uint32
	DTYPE_UINT32 = 0x03

	// Data type date/time
	DTYPE_DATETIME = 0x04

	// Data type float
	DTYPE_FLOAT = 0x05

	// Data type 128String
	// char string with 128 bytes, must be 0 terminated
	DTYPE_128STRING = 0x06
	// max char length = 127 + 0 termination
	STRING_PACKET_MAX_BUFFER_LENGTH = 127

	// Data type timestamp
	// in secondes since device was booted up, 32 bit unsigned integer (4 bytes)
	DTYPE_TIMESTAMP = 0x07

	// Data type 16bytes
	// 16 bytes of raw binary data
	DTYPE_16BYTES = 0x09
	// max byte length = 16 bytes
	BYTES16_PACKET_MAX_BUFFER_LENGTH = 16

	// Data type 66bytes
	// 66 bytes of raw binary data
	DTYPE_66BYTES = 0x08
	// max byte 65 bytes DATA TYPE name is a typo since this type is
	// used by the statemachine upload only, and that is 65 bytes
	BYTES66_PACKET_MAX_BUFFER_LENGTH_ = 65

	/* Error codes */

	// reserved: No error
	ERR_SUCCESS = 0x00

	// A request for an endpoint which does not exist on the device was received
	ERR_UNKNOWNEID = 0x01

	// WRITE was received for a readonly endpoint
	ERR_WRITEREADONLY = 0x02

	// A packet failed the CRC check
	// TODO How can we find out what information was lost?
	ERR_CRCFAILED = 0x03

	// A packet with a datatype that does not fit the endpoint was received
	ERR_DATATYPE = 0x04

	// A value was encountered that cannot be interpreted
	ERR_INVALID_VALUE = 0x05

	/* Networking */
	
	// hexabus default port
	PORT = "61616"

	// package transmit timeout
	NET_TIMEOUT = 3
	
)
