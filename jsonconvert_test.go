package jsonconvert

import (
	"github.com/howood/jsonconvert/internal/parser"
	"reflect"
	"testing"
)

type testData struct {
	Setting      string
	Input        string
	CheckData    string
	ResultHasErr bool
}

var responseTestData = map[string]testData{
	"test1": testData{
		ResultHasErr: false,
		Setting: `
{
	"GlossEntry": "$$glossary.GlossDiv.GlossList.GlossEntry",
	"GlossSeeAlso": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso",
	"GlossDef": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef",
	"title":  "$$glossary.title",
	 "key": "value"
}
`,
		Input: `
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
}`,

		CheckData: `
{
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
	},
	"GlossSeeAlso": ["GML", "XML"],
	"GlossDef": {
		"para": "A meta-markup language, used to create markup languages such as DocBook.",
		"GlossSeeAlso": ["GML", "XML"]
	},
	"title":  "example glossary",
	 "key": "value"
  }
`,
	},
	"test2": testData{
		ResultHasErr: false,
		Setting: `
{
	"billToaddress": ["$$[$$n].billTo.address"],
	"sku": ["$$[$$n].sku"],
	"key": "value"
}
`,
		Input: `
[
	{
		"billTo": {
			"address": "456 Oak Lanewwwww",
			"city": "Pretendvilledddd",
			"name": "Alice Brown33333",
			"state": "HI",
			"zip": "98999d"
		},
		"name": "Alice Brown33333",
		"price": 199.95,
		"shipTo": {
			"address": "456 Oak Lane",
			"city": "Pretendville",
			"name": "Bob Brown",
			"state": "HI",
			"zip": "98999"
		},
		"sku": "54321"
	},
	{
		"billTo": {
			"address": "456 Oak Lane",
			"city": "Pretendville",
			"name": "Alice Brown",
			"state": "HI",
			"zip": "98999"
		},
		"name": "Alice Brown",
		"price": 199.95,
		"shipTo": {
			"address": "456 Oak Lane",
			"city": "Pretendville",
			"name": "Bob Brown",
			"state": "HI",
			"zip": "98999"
		},
		"sku": "54321"
	},
	{
		"billTo": {
			"address": "123 Maple Street",
			"city": "Pretendville",
			"name": "John Smith",
			"state": "NY",
			"zip": "12345"
		},
		"name": "John Smith",
		"price": 23.95,
		"shipTo": {
			"address": "123 Maple Street",
			"city": "Pretendville",
			"name": "Jane Smith",
			"state": "NY",
			"zip": "12345"
		},
		"sku": "20223"
	}
]
`,
		CheckData: `
{
	"billToaddress": [
		"456 Oak Lanewwwww",
		"456 Oak Lane",
		"123 Maple Street"
	],
	"key": "value",
	"sku": [
		"54321",
		"54321",
		"20223"
	]
}
`,
	},
	"test3": testData{
		ResultHasErr: false,
		Setting: `
[
	{
		"$$recordset": "userdata",
		"$$joinrecordset": "orderdata",
		"$$joinrecordcolumn": "userid",
		"userid": "$$userid",
		"pref": "$$address.pref",
		"address": {
			"pref": "$$address.pref"
		},
		"name": "$$name",
		"orderdata": "$$joinrecordset"
	}
]
`,
		Input: `
{
	"userdata":
	[
		{
			"userid": 2,
			"address": {
				"pref": "tokyo"
			},
			"name": "aaa"
		},
		{
			"userid": 3,
			"address": {
				"pref": "osaka"
			},
			"name": "bbb"
		},
		{
			"userid": 4,
			"address": {
				"pref": "fukuoka"
			},
			"name": "ccc"
		}
	],
	"orderdata":
	[
		{
			"userid": 3,
			"orderid": 1,
			"product": "product_1"
		},
		{
			"userid": 4,
			"orderid": 2,
			"product": "product_2"
		},
		{
			"userid": 2,
			"orderid": 3,
			"product": "product_3"
		}
	]
}
`,
		CheckData: `
[
	{
		"address": {
			"pref": "tokyo"
		},
		"name": "aaa",
		"orderdata": {
			"orderid": 3,
			"product": "product_3",
			"userid": 2
		},
		"pref": "tokyo",
		"userid": 2
	},
	{
		"address": {
			"pref": "tokyo"
		},
		"name": "bbb",
		"orderdata": {
			"orderid": 1,
			"product": "product_1",
			"userid": 3
		},
		"pref": "osaka",
		"userid": 3
	},
	{
		"address": {
			"pref": "tokyo"
		},
		"name": "ccc",
		"orderdata": {
			"orderid": 2,
			"product": "product_2",
			"userid": 4
		},
		"pref": "fukuoka",
		"userid": 4
	}
]
`,
	},
	"test4": testData{
		ResultHasErr: false,
		Setting: `
{
	"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid",
			"userid": "$$userid",
			"pref": "$$address.pref",
			"address": {
				"pref": "$$address.pref"
			},
			"name": "$$name",
			"orderdata": "$$joinrecordset"
		}
	]
}
`,
		Input: `
{
	"userdata":
	[
		{
			"userid": 2,
			"address": {
				"pref": "tokyo"
			},
			"name": "aaa"
		},
		{
			"userid": 3,
			"address": {
				"pref": "osaka"
			},
			"name": "bbb"
		},
		{
			"userid": 4,
			"address": {
				"pref": "fukuoka"
			},
			"name": "ccc"
		}
	],
	"orderdata":
	[
		{
			"userid": 3,
			"orderid": 1,
			"product": "product_1"
		},
		{
			"userid": 4,
			"orderid": 2,
			"product": "product_2"
		},
		{
			"userid": 2,
			"orderid": 3,
			"product": "product_3"
		}
	]
}
`,
		CheckData: `
{
	"user": [
		{
			"address": {
				"pref": "tokyo"
			},
			"name": "aaa",
			"orderdata": {
				"orderid": 3,
				"product": "product_3",
				"userid": 2
			},
			"pref": "tokyo",
			"userid": 2
		},
		{
			"address": {
				"pref": "tokyo"
			},
			"name": "bbb",
			"orderdata": {
				"orderid": 1,
				"product": "product_1",
				"userid": 3
			},
			"pref": "osaka",
			"userid": 3
		},
		{
			"address": {
				"pref": "tokyo"
			},
			"name": "ccc",
			"orderdata": {
				"orderid": 2,
				"product": "product_2",
				"userid": 4
			},
			"pref": "fukuoka",
			"userid": 4
		}
	]
}
`,
	},
	"test5": testData{
		ResultHasErr: false,
		Setting: `
[
	{
		"$$recordset": "userdata",
		"userid": "$$userid",
		"name": "$$name"
	}
]
`,
		Input: `
{
	"userdata":
	[
		{
			"userid": 2,
			"address": {
				"pref": "tokyo"
			},
			"name": "aaa"
		},
		{
			"userid": 3,
			"address": {
				"pref": "osaka"
			},
			"name": "bbb"
		},
		{
			"userid": 4,
			"address": {
				"pref": "fukuoka"
			},
			"name": "ccc"
		}
	],
	"orderdata":
	[
		{
			"userid": 3,
			"orderid": 1,
			"product": "product_1"
		},
		{
			"userid": 4,
			"orderid": 2,
			"product": "product_2"
		},
		{
			"userid": 2,
			"orderid": 3,
			"product": "product_3"
		}
	]
}
`,
		CheckData: `
[
	{
		"name": "aaa",
		"userid": 2
	},
	{
		"name": "bbb",
		"userid": 3
	},
	{
		"name": "ccc",
		"userid": 4
	}
]
`,
	},
	"test6": testData{
		ResultHasErr: true,
		Setting:      "wwwww",
		Input:        "wwwww",
		CheckData:    "wwwww",
	},
	"test7": testData{
		ResultHasErr: true,
		Setting:      `{"aaa":7}`,
		Input:        "wwwww",
		CheckData:    "wwwww",
	},
	"test8": testData{
		ResultHasErr: true,
		Setting:      `{"data": "$$aaa.vvv"}`,
		Input:        `{"aaa":7}`,
		CheckData:    `{"aaa":7}`,
	},
	"test9": testData{
		ResultHasErr: true,
		Setting:      `[{"data": "$$aaa.vvv"}]`,
		Input:        `{"aaa":7}`,
		CheckData:    `{"aaa":7}`,
	},
	"test10": testData{
		ResultHasErr: false,
		Setting:      `[{"data": "$$aaa"}]`,
		Input:        `{"aaa":7}`,
		CheckData:    `[{"data":7}]`,
	},
	"test11": testData{
		ResultHasErr: true,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid"
		}
		]
}`,
		Input: `{"data": {"aaa":7}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
	"test12": testData{
		ResultHasErr: false,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid"
		}
		]
}`,
		Input: `{"userdata": {"aaa":7}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
	"test13": testData{
		ResultHasErr: false,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid"
		}
		]
}`,
		Input: `{"userdata": {"aaa":7}, "orderdata":{"aaa":7}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
	"test14": testData{
		ResultHasErr: true,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata"
		}
		]
}`,
		Input: `{"userdata": {"aaa":7}, "orderdata":{"aaa":7}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
	"test15": testData{
		ResultHasErr: false,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid"
		}
		]
}`,
		Input: `{"userdata": {"aaa":7}, "orderdata":{}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
	"test16": testData{
		ResultHasErr: false,
		Setting: `{
			"billToaddress": ["$$ccc.[$$n].billTo"]
		}
		`,
		Input: `{"ccc": [{"billTo":8},{"billTo":9}]}`,
		CheckData: `{
            "billToaddress": [
                8,
                9
            ]
        }`,
	},
	"test17": testData{
		ResultHasErr: true,
		Setting: `{
			"billToaddress": ["$$ccc2.[$$n].billTo"]
		}
		`,
		Input: `{"ccc": [{"billTo":8},{"billTo":9}]}`,
		CheckData: `{
            "billToaddress": [
                8,
                9
            ]
        }`,
	},
	"test18": testData{
		ResultHasErr: true,
		Setting: `{
			"billToaddress": ["$$ccc.[$$n].billTo2"]
		}
		`,
		Input: `{"ccc": [{"billTo":8},{"billTo":9}]}`,
		CheckData: `{
            "billToaddress": [
                8,
                9
            ]
        }`,
	},
	"test19": testData{
		ResultHasErr: false,
		Setting: `{
			"billToaddress": ["$$ccc.[$$n].billTo"]
		}
		`,
		Input: `{"ccc": {"billTo":8}}`,
		CheckData: `{
            "billToaddress": [
                8
            ]
        }`,
	},
	"test20": testData{
		ResultHasErr: true,
		Setting: `{
			"billToaddress": ["$$ccc.[$$n].billTo"]
		}
		`,
		Input:     `{"ccc": ""}`,
		CheckData: ``,
	},
	"test21": testData{
		ResultHasErr: false,
		Setting: `{"user": [
		{
			"$$recordset": "userdata",
			"$$joinrecordset": "orderdata",
			"$$joinrecordcolumn": "userid"
		}
		]
}`,
		Input: `{"userdata": {"userid":8}, "orderdata":{"userid": [8,9], "ssss": [4,5,7]}}`,
		CheckData: `{
            "user": [
                {}
            ]
        }`,
	},
}

func Test_JsonConvertTest1(t *testing.T) {
	jc := NewJSONConvert()
	for k, v := range responseTestData {
		jc.SetResponseSetting(k, v.Setting)
	}
	for k, v := range responseTestData {
		resultbyte, err := jc.Convert([]byte(v.Input), k)
		if (err != nil) != v.ResultHasErr {
			t.Fatalf("failed test :%s %#v", k, err)
		} else {
			t.Logf("failed test :%s %#v", k, err)
		}
		if v.ResultHasErr == false && checkEqualJSONByte(resultbyte, []byte(v.CheckData), t) == false {
			t.Log(string(resultbyte))
			t.Log(v.CheckData)
			t.Fatalf("failed JsonConvert :%s", k)
		}
		t.Logf("success : %s", k)

	}
	t.Log("success JsonConvert")

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
