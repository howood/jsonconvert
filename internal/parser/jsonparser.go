package parser

import "encoding/json"

const (
	// marshalPrefix is prefix for indented line beginning with
	marshalPrefix = ""
	// marshalPrefix is indented according to the indentation nesting
	marshalIndent = "    "
)

// JSONToByte is convert json struct to bytes
func JSONToByte(jsondata interface{}) ([]byte, error) {
	return json.MarshalIndent(jsondata, marshalPrefix, marshalIndent)
}

// ByteToJSON is convert bytes to json interface{}
func ByteToJSON(jsonbyte []byte) (interface{}, error) {
	var jsondata interface{}
	err := json.Unmarshal(jsonbyte, &jsondata)
	return jsondata, err
}

// ByteToJSONStruct is convert bytes to struct
func ByteToJSONStruct(jsonbyte []byte, jsonobj interface{}) error {
	return json.Unmarshal(jsonbyte, &jsonobj)
}
