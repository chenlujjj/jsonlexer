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
		input: `{"foo":  {"bar": "baz"}  }`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemLeftBrace, "{"}, {itemString, `"bar"`}, {itemColon, ":"}, {itemString, `"baz"`}, {itemRightBrace, "}"}, {itemRightBrace, "}"}},
	},
	{
		input: `["foo", true, false, null]`,
		items: []item{{itemLeftBracket, "["}, {itemString, `"foo"`}, {itemComma, ","}, {itemTrue, "true"}, {itemComma, ","}, {itemFalse, "false"}, {itemComma, ","}, {itemNull, "null"},{itemRightBracket, "]"}},
	},
	{
		input: `{"foo" :  ["bar", true]}`,
		items: []item{{itemLeftBrace, "{"}, {itemString, `"foo"`}, {itemColon, ":"}, {itemLeftBracket, "["}, {itemString, `"bar"`}, {itemComma, ","}, {itemTrue, "true"},{itemRightBracket, "]"}, {itemRightBrace, "}"}},
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
