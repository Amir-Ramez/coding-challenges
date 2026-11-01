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
		{"empty object", "{}", false},
		{"empty object with whitespace", " { } ", false},
		{"missing closing brace", "{", true},
		{"missing opening brace", "}", true},
		{"empty input", "", true},
		{"extra closing brace", "{}}", true},
		{"two objects", "{}{}", true},
		{"invalid character", "{}x", true},
		{"array", "[]", true},
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
