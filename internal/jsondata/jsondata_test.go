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

	json, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	jsonbyte, err := parser.JSONToByte(json)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if checkEqualJsonByte(jsonbyte, []byte(jsondatatestcheck), t) == false {
		t.Fatalf("failed JsonConvert ")
	}

	if _, err := jd.Query("glossary.GlossDiv.GlossList.GlossEntry2"); err == nil {
		t.Fatal("failed test no error")
	} else {
		t.Logf("failed test %#v", err)
	}

	t.Log("success JsonData")
}

func checkEqualJsonByte(input1, input2 []byte, t *testing.T) bool {
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
