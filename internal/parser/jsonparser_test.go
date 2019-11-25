package parser

import (
	"reflect"
	"testing"
)

func Test_JSONToByte_ByteToJSON(t *testing.T) {
	inputdata1 := map[string]interface{}{
		"message1": "ok",
		"message2": "ng",
	}
	result, err := JSONToByte(inputdata1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log("success JSONToByte")
	t.Log(string(result))
	resultmap, err := ByteToJSON(result)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(resultmap)
	if reflect.DeepEqual(inputdata1, resultmap) == false {
		t.Fatalf("failed ByteToJSON ")
	}
	t.Log("success ByteToJSON")
}

func Test_JSONToByte_ByteToJSON2(t *testing.T) {
	var inputdata1 map[string]interface{}
	result, err := JSONToByte(inputdata1)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log("success JSONToByte")
	t.Log(string(result))
	resultmap, err := ByteToJSON(result)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	t.Log(resultmap)
	if resultmap != nil {
		t.Fatalf("failed ByteToJSON ")
	}
	t.Log("success ByteToJSON")
}
