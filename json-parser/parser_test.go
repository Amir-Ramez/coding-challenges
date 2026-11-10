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
