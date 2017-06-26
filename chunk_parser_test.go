package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

var chunkTests = []struct{ in, expected string }{
	{"{{12}}", "12"},
	{"{{x}}", "123"},
	{"{{page.title}}", "Introduction"},
	{"{{ar[1]}}", "second"},
	{"{%if true%}true{%endif%}", "true"},
	{"{%if false%}false{%endif%}", ""},
	{"{%if 0%}true{%endif%}", "true"},
	{"{%if 1%}true{%endif%}", "true"},
	{"{%if x%}true{%endif%}", "true"},
	{"{%if y%}true{%endif%}", ""},
	{"{%if true%}true{%endif%}", "true"},
	{"{%if false%}false{%endif%}", ""},
	{"{%if true%}true{%else%}false{%endif%}", "true"},
	{"{%if false%}false{%else%}true{%endif%}", "true"},
	{"{%if true%}0{%elsif true%}1{%else%}2{%endif%}", "0"},
	{"{%if false%}0{%elsif true%}1{%else%}2{%endif%}", "1"},
	{"{%if false%}0{%elsif false%}1{%else%}2{%endif%}", "2"},
	{"{%unless true%}false{%endif%}", ""},
	{"{%unless false%}true{%endif%}", "true"},
	{"{%unless true%}false{%else%}true{%endif%}", "true"},
	{"{%unless false%}true{%else%}false{%endif%}", "true"},
	{"{%unless false%}0{%elsif true%}1{%else%}2{%endif%}", "0"},
	{"{%unless true%}0{%elsif true%}1{%else%}2{%endif%}", "1"},
	{"{%unless true%}0{%elsif false%}1{%else%}2{%endif%}", "2"},
	// {"{%for a in ar%}{{a}} {{%endfor%}", "first second third "},
}

var chunkTestContext = Context{map[string]interface{}{
	"x":  123,
	"ar": []string{"first", "second", "third"},
	"page": map[string]interface{}{
		"title": "Introduction",
	},
},
}

func TestChunkParser(t *testing.T) {
	for i, test := range chunkTests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			tokens := ScanChunks(test.in, "")
			// fmt.Println(tokens)
			ast, err := Parse(tokens)
			require.NoErrorf(t, err, test.in)
			// fmt.Println(MustYAML(ast))
			buf := new(bytes.Buffer)
			err = ast.Render(buf, chunkTestContext)
			require.NoErrorf(t, err, test.in)
			require.Equalf(t, test.expected, buf.String(), test.in)
		})
	}
}
