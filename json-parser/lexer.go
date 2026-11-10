package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrLexerAborted = errors.New("lexer aborted due to an error")

type TokenType int

const (
	LBrace TokenType = iota
	RBrace
	String
	Colon
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

func (l *Lexer) errorToken(err error) Token {
	l.err = ErrLexerAborted
	return Token{tokenType: Error, text: err.Error()}
}

func (l *Lexer) readString() Token {
	var b strings.Builder
	b.WriteRune('"')

	l.advance()

	for {
		if l.err != nil {
			if errors.Is(l.err, io.EOF) {
				return l.errorToken(fmt.Errorf("unterminated string"))
			}

			return l.errorToken(l.err)
		}

		c := l.curChar

		if c == '"' {
			b.WriteRune(c)
			l.advance()
			return Token{tokenType: String, text: b.String()}
		}

		// character from unicode code 0 through 30 are control characters
		if c < 0x20 {
			return l.errorToken(fmt.Errorf("unescaped control character in string"))
		}

		b.WriteRune(c)
		l.advance()

		if c == '\\' {
			if l.err != nil {
				return l.errorToken(fmt.Errorf("unterminated escape sequence"))
			}

			escaped := l.curChar

			switch escaped {
			case '"', '\\', '/', 'f', 'r', 't', 'n', 'b':
				b.WriteRune(escaped)
				l.advance()

			case 'u':
				b.WriteRune(escaped)
				l.advance()

				for range 4 {
					if l.err != nil || !isHexDigit(l.curChar) {
						return l.errorToken(fmt.Errorf("invalid Unicode character"))
					}

					b.WriteRune(l.curChar)
					l.advance()
				}

			default:
				return l.errorToken(fmt.Errorf("invalid escape sequence"))
			}
		}
	}
}

func isHexDigit(c rune) bool {
	return (c >= '0' && c <= '9') ||
		(c >= 'A' && c <= 'F') ||
		(c >= 'a' && c <= 'f')
}

func (l *Lexer) getNextToken() Token {
	l.skipWhiteSpace()

	if l.err != nil {
		if errors.Is(l.err, io.EOF) {
			return Token{tokenType: EOF}
		}

		return l.errorToken(l.err)
	}

	c := l.curChar

	switch c {
	case '{':
		l.advance()
		return Token{tokenType: LBrace, text: "{"}

	case '}':
		l.advance()
		return Token{tokenType: RBrace, text: "}"}

	case ':':
		l.advance()
		return Token{tokenType: Colon, text: ":"}
	case '"':
		return l.readString()
	default:
		l.advance()
		return Token{tokenType: Error, text: string(c)}
	}
}
