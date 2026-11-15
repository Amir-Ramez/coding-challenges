package main

import (
	"fmt"
	"io"
)

type Parser struct {
	CurrentToken Token
	Lexer        *Lexer
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		Lexer: NewLexer(r),
	}
}

func (p *Parser) advance() {
	p.CurrentToken = p.Lexer.getNextToken()
}

func (p *Parser) expect(want TokenType) error {
	if p.CurrentToken.tokenType == Error {
		return fmt.Errorf("lexer error: %s", p.CurrentToken.text)
	}

	if p.CurrentToken.tokenType != want {
		return fmt.Errorf("expected token %v, got %v", want, p.CurrentToken.tokenType)
	}

	p.advance()
	return nil
}

func (p *Parser) parseValue() error {
	switch p.CurrentToken.tokenType {
	case String:
		return p.expect(String)
	case Numeric:
		return p.expect(Numeric)
	case True:
		return p.expect(True)
	case False:
		return p.expect(False)
	case Null:
		return p.expect(Null)
	case LBrace:
		return p.parseObject()
	case LBracket:
		return p.parseArray()
	default:
		return fmt.Errorf("expected value, got %v", p.CurrentToken.tokenType)
	}
}

func (p *Parser) parseObject() error {
	if err := p.expect(LBrace); err != nil {
		return err
	}

	// empty object
	if p.CurrentToken.tokenType == RBrace {
		p.advance()
		return nil
	}

	for {
		if err := p.expect(String); err != nil {
			return err
		}

		if err := p.expect(Colon); err != nil {
			return err
		}

		if err := p.parseValue(); err != nil {
			return err
		}

		if err := p.expect(Comma); err != nil {
			break
		}
	}

	if err := p.expect(RBrace); err != nil {
		return err
	}

	return nil
}

func (p *Parser) parseArray() error {
	if err := p.expect(LBracket); err != nil {
		return err
	}

	if p.CurrentToken.tokenType == RBracket {
		p.advance()
		return nil
	}

	for {
		if err := p.parseValue(); err != nil {
			return err
		}

		if err := p.expect(Comma); err != nil {
			break
		}
	}

	if err := p.expect(RBracket); err != nil {
		return err
	}

	return nil
}

func (p *Parser) Parse() error {
	p.advance()

	if err := p.parseObject(); err != nil {
		return err
	}

	if p.CurrentToken.tokenType != EOF {
		return fmt.Errorf("unexpected token after JSON object: %q", p.CurrentToken.text)
	}

	return nil
}
