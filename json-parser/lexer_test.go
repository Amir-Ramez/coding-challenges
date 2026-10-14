package main

import (
	"strings"
	"testing"
)

func TestGetNextToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []Token
	}{
		{
			name:  "empty object",
			input: "{}",
			want: []Token{
				{tokenType: LBrace, text: "{"},
				{tokenType: RBrace, text: "}"},
				{tokenType: EOF},
			},
		},

		{
			name:  "object with whitespaces",
			input: " {\t\n} ",
			want: []Token{
				{tokenType: LBrace, text: "{"},
				{tokenType: RBrace, text: "}"},
				{tokenType: EOF},
			},
		},

		{
			name:  "empty input",
			input: "",
			want: []Token{
				{tokenType: EOF},
			},
		},

		{
			name:  "invalid character",
			input: "x",
			want: []Token{
				{tokenType: Error, text: "x"},
				{tokenType: EOF},
			},
		},

		{
			name:  "invalid whitespace",
			input: "\v",
			want: []Token{
				{tokenType: Error, text: "\v"},
				{tokenType: EOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(strings.NewReader(tt.input))

			for i, want := range tt.want {
				got := l.getNextToken()

				if got.tokenType != want.tokenType {
					t.Errorf("token %d: type = %v, want %v", i, got.tokenType, want.tokenType)
				}

				if got.text != want.text {
					t.Errorf("token %d: text = %v, want %v", i, got.text, want.text)
				}
			}
		})
	}
}
