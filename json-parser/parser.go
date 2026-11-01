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

func (p *Parser) parseObject() error {
	if err := p.expect(LBrace); err != nil {
		return err
	}

	if err := p.expect(RBrace); err != nil {
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
