package main

import (
	"bufio"
	"errors"
	"io"
)

type TokenType int

const (
	LBrace TokenType = iota
	RBrace
	EOF
	Error
)

type Token struct {
	tokenType TokenType
	text      string
}

type Lexer struct {
	r       *bufio.Reader
	curChar rune
	err     error
}

func NewLexer(r io.Reader) *Lexer {
	l := &Lexer{r: bufio.NewReader(r)}
	l.advance()
	return l
}

func (l *Lexer) advance() {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		l.curChar = 0
		l.err = err
		return
	}

	l.curChar = ch
	l.err = nil
}

func isJSONWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) skipWhiteSpace() {
	for l.err == nil && isJSONWhitespace(l.curChar) {
		l.advance()
	}
}

func (l *Lexer) getNextToken() Token {
	l.skipWhiteSpace()

	if l.err != nil {
		if errors.Is(l.err, io.EOF) {
			return Token{tokenType: EOF}
		}

		err := l.err
		l.err = io.EOF
		return Token{tokenType: Error, text: err.Error()}
	}

	c := l.curChar
	l.advance()

	switch c {
	case '{':
		return Token{tokenType: LBrace, text: "{"}
	case '}':
		return Token{tokenType: RBrace, text: "}"}
	default:
		return Token{tokenType: Error, text: string(c)}
	}
}
