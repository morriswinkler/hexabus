package hexabus

import "testing"
import "bytes"

type payload struct {
	data interface{}
}

var data_types = map[byte]payload{
	DTYPE_BOOL: payload{
		true,
	},
	DTYPE_UINT8: payload{
		uint8(109),
	},
	DTYPE_UINT32: payload{
		uint32(32434353),
	},
	DTYPE_DATETIME: payload{
		DateTime{17, 2, 15, 6, 3, 2014, 4},
	},
	DTYPE_FLOAT: payload{
		float32(10.102930),
	},
	DTYPE_128STRING: payload{
		"this is a hexabus packet test",
	},
	DTYPE_TIMESTAMP: payload{
		Timestamp{899992},
	},
	DTYPE_16BYTES: payload{
		make_byte_slice(16),
	},
	DTYPE_66BYTES: payload{
		make_byte_slice(65),
	},
}

var error_t = []byte{ERR_SUCCESS, ERR_UNKNOWNEID, ERR_WRITEREADONLY, ERR_CRCFAILED, ERR_DATATYPE, ERR_INVALID_VALUE}

func make_byte_slice(size int) []byte {
	byte_slice := make([]byte, size)
	for i := 0; i < size; i++ {
		byte_slice[i] = byte(i)
	}
	return byte_slice
}

func Test_ErrorPacket(t *testing.T) {

	for _, v := range error_t {
		p_error := ErrorPacket{FLAG_NONE, v}
		packet := p_error.Encode()

		p0_error := ErrorPacket{}
		p0_error.Decode(packet)

		if p0_error != p_error {
			t.Errorf("ErrorPacket with error Type %x did not match while testing: \n Encode: %+v \n Decode: %+v \n", v, p_error, p0_error)
		} else {
			t.Logf("ErrorPacket with Err type %x passed test", v)
			t.Logf("Send    :%+v", p_error)
			t.Logf("Receive :%+v", p0_error)
			t.Logf("RAW     :%x", packet)
			t.Logf("")
			
		}
	}
}

func Test_InfoPacket(t *testing.T) {
	for k, v := range data_types {
		p_info := InfoPacket{FLAG_NONE, 10, k, v.data}
		packet, err := p_info.Encode()
		if err != nil {
			t.Errorf("%s", err)
		}

		p0_info := InfoPacket{}
		p0_info.Decode(packet)

		if k != DTYPE_16BYTES && k != DTYPE_66BYTES {
			if  p_info != p0_info {
				t.Errorf("InfoPacket with Data type %d did not match while testing: \n Encode: %+v \n Decode: %+v \n", p0_info.Dtype, p_info, p0_info)
			} else {
				t.Logf("InfoPacket with Data type %d passed test", k)
				t.Logf("Send    :%+v", p_info)
				t.Logf("Receive :%+v", p0_info)
				t.Logf("RAW     :%x", packet)
				t.Logf("")
			}
		} else if k == DTYPE_16BYTES || k == DTYPE_66BYTES {
			if ( bytes.Equal(p_info.Data.([]byte), p0_info.Data.([]byte)) == false) {
				t.Errorf("InfoPacket with Data type %d did not match while testing: \n Encode: %+v \n Decode: %+v \n", p0_info.Dtype, p_info, p0_info)
			} else {
				t.Logf("InfoPacket with Data type %d passed test", k)
				t.Logf("Send    :%+v", p_info)
				t.Logf("Receive :%+v", p0_info)
				t.Logf("RAW     :%x", packet)
				t.Logf("")

			}
		}
	}
}

var eids = []uint32{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16}

func Test_QueryPacket(t *testing.T) {

	
	for _, v := range eids {
		p_query := QueryPacket{FLAG_NONE, v}
		packet := p_query.Encode()

		p0_query := QueryPacket{}
		p0_query.Decode(packet)

		if p0_query != p_query {
			t.Errorf("QueryPacket with EID %d did not match while testing: \n Encode: %+v \n Decode: %+v \n", v, p_query, p0_query)
		} else {
			t.Logf("QuerryPacket with EID type %d passed test", v)
			t.Logf("Send    :%+v", p_query)
			t.Logf("Receive :%+v", p0_query)
			t.Logf("RAW     :%x", packet)
			t.Logf("")
			
		}
	}
}
