package jsonconvert

import "encoding/json"

const MARSHAL_PREFIX = ""
const MARSHAL_INDENT = "    "

func jsonToByte(jsondata interface{}) ([]byte, error) {
	return json.MarshalIndent(jsondata, MARSHAL_PREFIX, MARSHAL_INDENT)
}

func byteToJson(jsonbyte []byte) (interface{}, error) {
	var jsondata interface{}
	err := json.Unmarshal(jsonbyte, &jsondata)
	return jsondata, err
}

func byteToJsonStruct(jsonbyte []byte, jsonobj interface{}) error {
	return json.Unmarshal(jsonbyte, &jsonobj)
}
