package main

import (
	"strings"
	"testing"
)

func TestParseObject(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Empty objects
		{"empty object", "{}", false},
		{"empty object with whitespace", " { } ", false},

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
		{"multiple string pairs with whitespace", ` { "name" : "Amir" , "salary" : "2000" } `, false},

		// Boolean values
		{"true value", `{"active":true}`, false},
		{"false value", `{"active":false}`, false},
		{"true value with whitespace", `{"active" : true}`, false},
		{"false value with whitespace", `{"active" : false}`, false},

		// Null values
		{"null value", `{"value":null}`, false},
		{"null value with whitespace", `{"value" : null}`, false},

		// Numeric values: integers
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

		// Missing / malformed object structure
		{"missing closing brace", "{", true},
		{"missing opening brace", "}", true},
		{"empty input", "", true},
		{"extra closing brace", "{}}", true},
		{"two objects", "{}{}", true},
		{"invalid character", "{}x", true},
		{"array", "[]", true},

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			err := p.Parse()

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
