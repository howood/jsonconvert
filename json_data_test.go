package jsonconvert

import (
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
	jd, err := newJSONData(josnDataTest)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}

	title, err := jd.query("glossary.title")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if title != "example glossary" {
		t.Fatalf("failed JsonData: get string")
	}

	array, err := jd.query("glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso")
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	checkarray := []interface{}{"GML", "XML"}
	if reflect.DeepEqual(array, checkarray) == false {
		t.Fatalf("failed JsonData ")
	}

	json, err := jd.query("glossary.GlossDiv.GlossList.GlossEntry")
	jsonbyte, err := jsonToByte(json)
	if err != nil {
		t.Fatalf("failed test %#v", err)
	}
	if checkEqualJsonByte(jsonbyte, []byte(jsondatatestcheck), t) == false {
		t.Fatalf("failed JsonConvert ")
	}

	t.Log("success JsonData")
}
