package main

import (
	"strings"
	"testing"
)

func runParseTest(t *testing.T, input string, wantErr bool) {
	t.Helper()

	p := NewParser(strings.NewReader(input))
	err := p.Parse()

	if (err != nil) != wantErr {
		t.Errorf("Parse() error = %v, wantErr %v", err, wantErr)
	}
}

func TestParseObject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Empty objects
		{"empty object", `{}`, false},
		{"empty object with whitespace", ` { } `, false},

		// String key/value pairs
		{"string key and value", `{"name":"Amir"}`, false},
		{"string key and value with whitespace", ` { "name" : "Amir" } `, false},
		{"empty string key and value", `{"":""}`, false},
		{"string with spaces", `{"name":"Amir Ramez"}`, false},
		{"string with escaped quote", `{"message":"hello \"world\""}`, false},
		{"string with escaped backslash", `{"path":"hello\\world"}`, false},
		{"string with unicode", `{"emoji":"😀"}`, false},

		// Multiple key/value pairs
		{"two string pairs", `{"name":"Amir","salary":"2000"}`, false},
		{"multiple pairs with whitespace", ` { "name" : "Amir" , "salary" : "2000" } `, false},

		// Boolean values
		{"true value", `{"active":true}`, false},
		{"false value", `{"active":false}`, false},
		{"true value with whitespace", `{"active" : true}`, false},
		{"false value with whitespace", `{"active" : false}`, false},

		// Null values
		{"null value", `{"value":null}`, false},
		{"null value with whitespace", `{"value" : null}`, false},

		// Numeric values
		{"zero", `{"count":0}`, false},
		{"positive integer", `{"count":42}`, false},
		{"negative integer", `{"count":-42}`, false},
		{"negative zero", `{"count":-0}`, false},
		{"large integer", `{"count":123456789}`, false},
		{"integer with whitespace", `{"count" : -42 }`, false},
		{"multiple numeric values", `{"min":1,"max":100}`, false},

		// Mixed value types
		{"string and boolean", `{"name":"Amir","active":true}`, false},
		{"boolean and string", `{"active":false,"name":"Amir"}`, false},
		{"string and null", `{"name":"Amir","value":null}`, false},
		{"boolean and null", `{"active":true,"value":null}`, false},
		{"all supported value types", `{"name":"Amir","active":true,"verified":false,"value":null}`, false},

		// Nested objects
		{"nested empty object", `{"profile":{}}`, false},
		{"nested object", `{"profile":{"active":true}}`, false},
		{"nested object with multiple pairs", `{"profile":{"name":"Amir","active":true}}`, false},
		{"multiple nested objects", `{"user":{"profile":{"active":true}}}`, false},
		{"nested object with numeric value", `{"settings":{"count":42}}`, false},
		{"nested object with null value", `{"data":{"value":null}}`, false},
		{"nested object with mixed values", `{"data":{"name":"Amir","count":42,"active":true,"value":null}}`, false},

		// Missing / malformed object structure
		{"missing closing brace", `{`, true},
		{"missing opening brace", `}`, true},
		{"empty input", ``, true},
		{"extra closing brace", `{}}`, true},
		{"two objects", `{}{}`, true},
		{"invalid character", `{}x`, true},
		{"array", `[]`, true},

		// Invalid string pairs
		{"missing key", `{: "value"}`, true},
		{"missing colon", `{"key" "value"}`, true},
		{"missing value", `{"key":}`, true},
		{"missing closing quote in key", `{"key:"value"}`, true},
		{"missing closing quote in value", `{"key":"value}`, true},
		{"unquoted key", `{key:"value"}`, true},
		{"unquoted value", `{"key":value}`, true},
		{"invalid escape", `{"key":"hello\q"}`, true},
		{"unterminated string", `{"key":"hello}`, true},
		{"unescaped control character", "{\"key\":\"hello\nworld\"}", true},

		// Invalid boolean/null values
		{"invalid boolean", `{"active":yes}`, true},
		{"invalid null", `{"value":nil}`, true},
		{"boolean missing value", `{"active":}`, true},
		{"null missing value", `{"value":}`, true},

		// Invalid numeric values
		{"minus without digits", `{"count":-}`, true},
		{"minus followed by non-digit", `{"count":-x}`, true},
		{"leading zero", `{"count":01}`, true},
		{"negative leading zero", `{"count":-01}`, true},

		// Invalid commas
		{"trailing comma", `{"name":"Amir",}`, true},
		{"missing comma", `{"name":"Amir" "salary":"2000"}`, true},
		{"leading comma", `{,"name":"Amir"}`, true},
		{"double comma", `{"name":"Amir",,"salary":"2000"}`, true},

		// Invalid nested objects
		{"nested object missing closing brace", `{"profile":{"active":true}`, true},
		{"nested object missing colon", `{"profile":{"active" true}}`, true},
		{"nested object missing value", `{"profile":{"active":}}`, true},
		{"nested object trailing comma", `{"profile":{"active":true,}}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParseTest(t, tt.input, tt.wantErr)
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Empty arrays
		{"empty array", `{"items":[]}`, false},
		{"empty array with whitespace", `{"items":[ ]}`, false},

		// Single-element arrays
		{"single string", `{"items":["Go"]}`, false},
		{"single number", `{"items":[42]}`, false},
		{"single negative number", `{"items":[-42]}`, false},
		{"single boolean true", `{"items":[true]}`, false},
		{"single boolean false", `{"items":[false]}`, false},
		{"single null", `{"items":[null]}`, false},
		{"single empty string", `{"items":[""]}`, false},

		// Multiple elements
		{"multiple strings", `{"items":["Go","C","Rust"]}`, false},
		{"multiple numbers", `{"items":[1,2,3,4,5]}`, false},
		{"multiple booleans", `{"items":[true,false,true]}`, false},
		{"multiple nulls", `{"items":[null,null]}`, false},
		{"array with whitespace", `{"items":[ 1 , 2 , 3 ]}`, false},

		// Mixed element types
		{"mixed values", `{"items":["Go",42,true,false,null]}`, false},
		{"mixed values with whitespace", `{"items":[ "Go" , 42 , true , null ]}`, false},

		// Arrays alongside other object values
		{"array and string", `{"items":[1,2],"name":"Amir"}`, false},
		{"multiple arrays", `{"first":[1,2],"second":[3,4]}`, false},

		// Invalid array structure
		{"missing closing bracket", `{"items":[1,2,3}`, true},
		{"missing opening bracket", `{"items":1,2,3]}`, true},
		{"missing both brackets", `{"items":1,2,3}`, true},
		{"array missing value", `{"items":[ ]}`, false},
		{"array missing first value", `{"items":[,1]}`, true},
		{"array missing value between commas", `{"items":[1,,2]}`, true},
		{"trailing comma", `{"items":[1,2,]}`, true},
		{"missing comma", `{"items":[1 2]}`, true},
		{"extra closing bracket", `{"items":[1]]}`, true},
		{"extra opening bracket", `{"items":[[1]}`, true},
		{"object closing brace instead of bracket", `{"items":[1,2}`, true},
		{"invalid array value", `{"items":[hello]}`, true},
		{"empty element", `{"items":[1, ,2]}`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParseTest(t, tt.input, tt.wantErr)
		})
	}
}

func TestParseNestedStructures(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Arrays containing objects
		{"array with one object", `{"users":[{"name":"Amir"}]}`, false},
		{"array with multiple objects", `{"users":[{"name":"Amir"},{"name":"John"}]}`, false},
		{"array with empty object", `{"users":[{}]}`, false},

		// Objects containing arrays
		{"object with empty array", `{"user":{"items":[]}}`, false},
		{"object with array", `{"user":{"items":[1,2,3]}}`, false},

		// Nested arrays
		{"nested empty array", `{"items":[[]]}`, false},
		{"nested arrays", `{"items":[[1,2],[3,4]]}`, false},
		{"deeply nested arrays", `{"items":[[[1]]]}`, false},

		// Mixed nesting
		{"object inside array inside object", `{"data":[{"active":true}]}`, false},
		{"array inside object inside array", `{"data":[{"items":[1,2]}]}`, false},
		{"object and array mixed deeply", `{"data":{"users":[{"name":"Amir","scores":[10,20]}]}}`, false},
		{"multiple nesting levels", `{"a":[{"b":{"c":[1,2,3]}}]}`, false},

		// Invalid nested structures
		{"array with malformed object", `{"items":[{"name"}]}`, true},
		{"array with malformed nested array", `{"items":[[1,2,]]}`, true},
		{"nested array missing closing bracket", `{"items":[[1,2]}`, true},
		{"nested object missing closing brace", `{"items":[{"name":"Amir"}]`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runParseTest(t, tt.input, tt.wantErr)
		})
	}
}
