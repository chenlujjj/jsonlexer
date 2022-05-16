package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testcase struct {
	input string
	items []item
}

var testcases = []testcase{
	{
		input: "   \t\n\r",
		items: nil,
	},
	{
		input: `"hello"`,
		items: []item{{itemString, `"hello"`}},
	},
	{
		input: "true",
		items: []item{{itemTrue, "true"}},
	},
	{
		input: "false",
		items: []item{{itemFalse, "false"}},
	},
	{
		input: "null",
		items: []item{{itemNull, "null"}},
	},
	{
		input: "3",
		items: []item{{itemNumber, "3"}},
	},
	{
		input: "-3",
		items: []item{{itemNumber, "-3"}},
	},
	{
		input: "3e10",
		items: []item{{itemNumber, "3e10"}},
	},
	{
		input: "-3e10",
		items: []item{{itemNumber, "-3e10"}},
	},
	{
		input: "-3e+2",
		items: []item{{itemNumber, "-3e+2"}},
	},
	{
		input: "-3e-02",
		items: []item{{itemNumber, "-3e-02"}},
	},

	// {
	// 	input: "3.14",
	// 	items: []item{{itemNumber, "3"}},
	// },

	{
		input: "{}",
		items: []item{{itemLeftBrace, "{"}, {itemRightBrace, "}"}},
	},
	{
		input: "{",
		items: []item{{itemLeftBrace, "{"}, {itemError, "unexpected eof"}},
	},
	{
		input: "{  \n\r\t}",
		items: []item{{itemLeftBrace, "{"}, {itemRightBrace, "}"}},
	},
	{
		input: `{"foo":  "bar"}`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemString, `"bar"`}, {itemRightBrace, "}"}},
	},
	{
		input: `{"foo":  "bar",    "baz": "qux"}`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemString, `"bar"`}, {itemComma, ","}, {itemString, `"baz"`}, {itemColon, ":"}, {itemString, `"qux"`}, {itemRightBrace, "}"}},
	},
	{
		input: `{"foo":  {"bar": 3e+03}  }`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemLeftBrace, "{"}, {itemString, `"bar"`}, {itemColon, ":"}, {itemNumber, "3e+03"}, {itemRightBrace, "}"}, {itemRightBrace, "}"}},
	},
	{
		input: `["foo", true, false, null, 3]`,
		items: []item{{itemLeftBracket, "["}, {itemString, `"foo"`}, {itemComma, ","}, {itemTrue, "true"}, {itemComma, ","}, {itemFalse, "false"}, {itemComma, ","}, {itemNull, "null"}, {itemComma, ","}, {itemNumber, "3"}, {itemRightBracket, "]"}},
	},
	{
		input: `{"foo" :  ["bar", true]}`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemLeftBracket, "["}, {itemString, `"bar"`}, {itemComma, ","}, {itemTrue, "true"}, {itemRightBracket, "]"}, {itemRightBrace, "}"}},
	},
	{
		input: `["foo", {"bar": true}]`,
		items: []item{{itemLeftBracket, "["}, {itemString, `"foo"`}, {itemComma, ","}, {itemLeftBrace, "{"}, {itemString, `"bar"`}, {itemColon, ":"}, {itemTrue, "true"}, {itemRightBrace, "}"}, {itemRightBracket, "]"}},
	},
}

func TestLex(t *testing.T) {
	assert := assert.New(t)

	for _, tc := range testcases {
		lex := lexer{input: tc.input}
		lex.run()
		assert.Equalf(tc.items, lex.items, "")
	}
}
