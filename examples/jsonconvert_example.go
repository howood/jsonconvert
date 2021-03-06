package main

import (
	"log"

	"github.com/howood/jsonconvert"
)

const (
	responseSetting = "setting"
	responseInput   = "input"
)

var responseTestData = map[string]map[string]string{
	"test1": {
		responseSetting: `
{
	"GlossEntry": "$$glossary.GlossDiv.GlossList.GlossEntry",
	"GlossSeeAlso": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso",
	"GlossDef": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef",
	"title":  "$$glossary.title",
	 "key": "value"
  }
`,
		responseInput: `
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
	},
	"test2": {
		responseSetting: `
{
	"billToaddress": ["$$[$$n].billTo.address"],
	"sku": ["$$[$$n].sku"],
	"key": "value"
}
`,
		responseInput: `
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
	},
	"test3": {
		responseSetting: `
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
		responseInput: `
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
	},
	"test4": {
		responseSetting: `
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
		responseInput: `
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
	},
	"test5": {
		responseSetting: `
[
	{
		"$$recordset": "userdata",
		"userid": "$$userid",
		"name": "$$name"
	}
]
`,
		responseInput: `
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
	},
}

func main() {
	jc := jsonconvert.NewJSONConvert()
	for k, v := range responseTestData {
		jc.SetResponseSetting(k, v[responseSetting])
	}
	for k, v := range responseTestData {
		resultbyte, err := jc.Convert([]byte(v[responseInput]), k)
		if err != nil {
			log.Fatalf("failed  :%s %#v", k, err)
		}
		log.Println(string(resultbyte))
		log.Printf("success : %s", k)

	}
	log.Print("success JsonConvert")

}
