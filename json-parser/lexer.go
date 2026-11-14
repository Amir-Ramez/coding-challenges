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
	Numeric
	True
	False
	Null
	Colon
	Comma
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

func (l *Lexer) readKeyword(value string, tokenType TokenType) Token {
	for i := range len(value) {
		if l.curChar != rune(value[i]) {
			return l.errorToken(fmt.Errorf("invalid value"))
		}

		l.advance()
	}

	return Token{tokenType: tokenType, text: value}
}

func (l *Lexer) readNumber() Token {
	value := make([]byte, 0)

	if l.curChar == '-' {
		value = append(value, byte(l.curChar))
		l.advance()

		if l.curChar < '0' || l.curChar > '9' {
			return l.errorToken(fmt.Errorf("invalid number"))
		}
	}

	if l.curChar == '0' {
		value = append(value, byte(l.curChar))
		l.advance()

		if l.curChar >= '0' && l.curChar <= '9' {
			return l.errorToken(fmt.Errorf("leading zeros are not allowed"))
		}
	} else {
		for l.curChar >= '0' && l.curChar <= '9' {
			value = append(value, byte(l.curChar))
			l.advance()
		}
	}

	return Token{tokenType: Numeric, text: string(value)}
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

	switch {
	case c == '{':
		l.advance()
		return Token{tokenType: LBrace, text: "{"}

	case c == '}':
		l.advance()
		return Token{tokenType: RBrace, text: "}"}

	case c == ':':
		l.advance()
		return Token{tokenType: Colon, text: ":"}

	case c == '"':
		return l.readString()

	case c == ',':
		l.advance()
		return Token{tokenType: Comma, text: ","}

	case c == 't':
		return l.readKeyword("true", True)

	case c == 'f':
		return l.readKeyword("false", False)

	case c == 'n':
		return l.readKeyword("null", Null)

	case (c >= '0' && c <= '9') || c == '-':
		return l.readNumber()

	default:
		l.advance()
		return Token{tokenType: Error, text: string(c)}
	}
}
