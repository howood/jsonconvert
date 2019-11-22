package jsonconvert

var responseTestData = map[string]map[string]string{
	"test1": map[string]string{
		RESPONSE_SETTING: `
{
	"GlossEntry": "$$glossary.GlossDiv.GlossList.GlossEntry",
	"GlossSeeAlso": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef.GlossSeeAlso",
	"GlossDef": "$$glossary.GlossDiv.GlossList.GlossEntry.GlossDef",
	"title":  "$$glossary.title",
	 "key": "value"
  }
`,
		RESPONSE_INPUT: `
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

		RESPONSE_CHECKDATA: `
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
	"test2": map[string]string{
		RESPONSE_SETTING: `
{
	"billToaddress": ["$$[$$n].billTo.address"],
	"sku": ["$$[$$n].sku"],
	"key": "value"
}
`,
		RESPONSE_INPUT: `
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
		RESPONSE_CHECKDATA: `
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
	"test3": map[string]string{
		RESPONSE_SETTING: `
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
		RESPONSE_INPUT: `
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
		RESPONSE_CHECKDATA: `
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
	"test4": map[string]string{
		RESPONSE_SETTING: `
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
		RESPONSE_INPUT: `
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
		RESPONSE_CHECKDATA: `
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
	"test5": map[string]string{
		RESPONSE_SETTING: `
[
	{
		"$$recordset": "userdata",
		"userid": "$$userid",
		"name": "$$name"
	}
]
`,
		RESPONSE_INPUT: `
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
		RESPONSE_CHECKDATA: `
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
}
