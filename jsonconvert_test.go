package jsonconvert

import (
	"reflect"
	"testing"
)

func Test_JsonConvertTest1(t *testing.T) {
	jc := NewJSONConvert()
	for k, v := range responseTestData {
		jc.SetResponse(k, v[RESPONSE_SETTING])
	}
	for k, v := range responseTestData {
		resultbyte, err := jc.Convert([]byte(v[RESPONSE_INPUT]), k)
		if err != nil {
			t.Fatalf("failed test :%s %#v", k, err)
		}
		if checkEqualJsonByte(resultbyte, []byte(v[RESPONSE_CHECKDATA]), t) == false {
			t.Log(string(resultbyte))
			t.Log(v[RESPONSE_CHECKDATA])
			t.Fatalf("failed JsonConvert :%s", k)
		}
		t.Logf("success : %s", k)

	}
	t.Log("success JsonConvert")

}

func checkEqualJsonByte(input1, input2 []byte, t *testing.T) bool {
	var json1 interface{}
	var json2 interface{}

	if err := byteToJsonStruct(input1, &json1); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if err := byteToJsonStruct(input2, &json2); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	return reflect.DeepEqual(json1, json2)
}
