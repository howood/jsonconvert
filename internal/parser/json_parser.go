package parser

import "encoding/json"

const (
	// MARSHAL_PREFIX is prefix for indented line beginning with
	MARSHAL_PREFIX = ""
	// MARSHAL_PREFIX is indented according to the indentation nesting
	MARSHAL_INDENT = "    "
)

// jsonToByte is convert json struct to bytes
func JsonToByte(jsondata interface{}) ([]byte, error) {
	return json.MarshalIndent(jsondata, MARSHAL_PREFIX, MARSHAL_INDENT)
}

// byteToJson is convert bytes to json interface{}
func ByteToJson(jsonbyte []byte) (interface{}, error) {
	var jsondata interface{}
	err := json.Unmarshal(jsonbyte, &jsondata)
	return jsondata, err
}

// byteToJson is convert bytes to struct
func ByteToJsonStruct(jsonbyte []byte, jsonobj interface{}) error {
	return json.Unmarshal(jsonbyte, &jsonobj)
}
