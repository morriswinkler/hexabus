package hexabus

import "testing"

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
		DateTime{17,2,15,6,3,2014,4},
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
}

var error = []byte{ERR_SUCESS, ERR_UNKNOWNEID, ERR_WRITEREADONLY, ERR_CRCFAILED, ERR_DATATYPE, ERR_INVALID_VALUE}

func make_byte_slice(size int) []byte {
	byte_slice := make([]byte, size)
	for i := 0; i < size; i++ {
		byte_slice[i] = byte(i)
	}
    return byte_slice
}

func Test_ErrorPacket(t *testing.T) {

	for _, v := range error {
		p_error := ErrorPacket{FLAG_NONE,v}
		packet := p_error.Encode()
		
		p0_error := ErrorPacket{}
		p0_error.Decode(packet)

		if p0_error != p_error {
			t.Errorf("ErrorPacket did not match while testing: \n Encode: %+v \n Decode: %+v \n", p_error, p0_error)
		} else {
			t.Log("ErrorPacket test passed")
		}
	}
}

func Test_InfoPacket(t *testing.T) {
	for k, v := range data_types{
		p_info := InfoPacket{FLAG_NONE,10,k, v.data}
		packet := p_info.Encode()
		
		p0_info := InfoPacket{}
		p0_info.Decode(packet)
		
		t.Log(k)
		if k != 9 || k != 8 {
			if p_info != p0_info {
				t.Errorf("InfoPacket with datatype %d did not match while testing: \n Encode: %+v \n Decode: %+v \n", p0_info.Dtype, p_info, p0_info)
			} else {
				t.Log("InfoPackte test passed")
			}
		}
	}
}
