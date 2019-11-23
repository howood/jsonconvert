package parser

import (
	"reflect"
	"testing"
)

func Test_JsonToByte_ByteToJson(t *testing.T) {
	inputdata1 := map[string]interface{}{
		"message1": "ok",
		"message2": "ng",
	}
	result, err := JsonToByte(inputdata1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log("success JsonToByte")
	t.Log(string(result))
	resultmap, err := ByteToJson(result)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(resultmap)
	if reflect.DeepEqual(inputdata1, resultmap) == false {
		t.Fatalf("failed ByteToJson ")
	}
	t.Log("success ByteToJson")
}

func Test_JsonToByte_ByteToJson2(t *testing.T) {
	var inputdata1 map[string]interface{}
	result, err := JsonToByte(inputdata1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log("success JsonToByte")
	t.Log(string(result))
	resultmap, err := ByteToJson(result)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(resultmap)
	if resultmap != nil {
		t.Fatalf("failed ByteToJson ")
	}
	t.Log("success ByteToJson")
}
