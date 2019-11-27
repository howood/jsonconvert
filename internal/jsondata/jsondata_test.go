package jsondata

import (
	"github.com/howood/jsonconvert/internal/parser"
	"reflect"
	"testing"
)

var josnDataTest = `
{
	"glossary": {
		"title": "example glossary",
		"GlossDiv": {
			"title": "S",
			"GlossList": {
				"GlossEntry": {
					"ID": "SGML",
					"SortAs": "SGML",
					"GlossTerm": "Standard Generalized Markup Language",
					"Acronym": "SGML",
					"Abbrev": "ISO 8879:1986",
					"GlossDef": {
						"para": "A meta-markup language, used to create markup languages such as DocBook.",
						"GlossSeeAlso": ["GML", "XML"]
					},
					"GlossSee": "markup"
				}
			}
		}
	}
}
`
var josnDataTestInner = `
{
	"title": "example glossary",
	"GlossDiv": {
		"title": "S",
		"GlossList": {
			"GlossEntry": {
				"ID": "SGML",
				"SortAs": "SGML",
				"GlossTerm": "Standard Generalized Markup Language",
				"Acronym": "SGML",
				"Abbrev": "ISO 8879:1986",
				"GlossDef": {
					"para": "A meta-markup language, used to create markup languages such as DocBook.",
					"GlossSeeAlso": ["GML", "XML"]
				},
				"GlossSee": "markup"
			}
		}
	}
}
`

var jsondatatestcheck = `
{
	"ID": "SGML",
	"SortAs": "SGML",
	"GlossTerm": "Standard Generalized Markup Language",
	"Acronym": "SGML",
	"Abbrev": "ISO 8879:1986",
	"GlossDef": {
		"para": "A meta-markup language, used to create markup languages such as DocBook.",
		"GlossSeeAlso": ["GML", "XML"]
	},
	"GlossSee": "markup"
}
`

func Test_JsonData(t *testing.T) {
	var err error
	if _, err := NewJSONData(""); err == nil {
		t.Fatal("failed test string", err)
	} else {
		t.Logf("failed test %#v", err)
	}

	if _, err := NewJSONData([]byte("s")); err == nil {
		t.Fatal("failed test bytes", err)
	} else {
		t.Logf("failed test %#v", err)
	}

	jd, err := NewJSONData(josnDataTest)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	jsondata, err := jd.Query(".")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	var checkalldata interface{}
	if err := parser.ByteToJSONStruct([]byte(josnDataTest), &checkalldata); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if reflect.DeepEqual(checkalldata, jsondata) == false {
		t.Fatalf("failed JsonData: getall")
	}

	if _, err := jd.Query(""); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	if _, err := jd.Query("sssss"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	josndata2, err := jd.Query("glossary")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	var checkalldata2 interface{}
	if err := parser.ByteToJSONStruct([]byte(josnDataTestInner), &checkalldata2); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if reflect.DeepEqual(checkalldata2, josndata2) == false {
		t.Fatalf("failed JsonData: getall")
	}

	title, err := jd.Query("glossary.title")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if title != "example glossary" {
		t.Fatalf("failed JsonData: get string")
	}

	array, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	checkarray := []interface{}{"GML", "XML"}
	if reflect.DeepEqual(array, checkarray) == false {
		t.Fatalf("failed JsonData ")
	}

	arraydata, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.[0]")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if arraydata != "GML" {
		t.Fatalf("failed JsonData ")
	}

	if _, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.[2]"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	if _, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso.[aaa]"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	if _, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.[1]"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	json, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	jsonbyte, err := parser.JSONToByte(json)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if checkEqualJSONByte(jsonbyte, []byte(jsondatatestcheck), t) == false {
		t.Fatalf("failed JsonConvert ")
	}

	if _, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry2"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	if _, err := NewJSONData([]byte(josnDataTest)); err != nil {
		t.Fatalf("failed test %#v", err)
	}

	jsonmap := map[string]interface{}{
		"aaa": "aa",
		"bbb": "bb",
	}
	if _, err := NewJSONData(jsonmap); err != nil {
		t.Fatalf("failed test %#v", err)
	}

	t.Log("success JsonData")
}

func checkEqualJSONByte(input1, input2 []byte, t *testing.T) bool {
	var json1 interface{}
	var json2 interface{}

	if err := parser.ByteToJSONStruct(input1, &json1); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if err := parser.ByteToJSONStruct(input2, &json2); err != nil {
		t.Fatalf("failed test %#v", err)
	}
	return reflect.DeepEqual(json1, json2)
}
