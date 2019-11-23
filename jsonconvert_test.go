package jsonconvert

import (
	"reflect"
	"testing"
)

func Test_JsonConvertTest1(t *testing.T) {
	jc := NewJSONConvert()
	for k, v := range responseTestData {
		jc.SetResponseSetting(k, v["setting"])
	}
	for k, v := range responseTestData {
		resultbyte, err := jc.Convert([]byte(v["input"]), k)
		if err != nil {
			t.Fatalf("failed test :%s %#v", k, err)
		}
		if checkEqualJsonByte(resultbyte, []byte(v["checkdata"]), t) == false {
			t.Log(string(resultbyte))
			t.Log(v["checkdata"])
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
