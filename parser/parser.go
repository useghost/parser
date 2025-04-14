package parser

import (
	"fmt"
	"ghostlang/ast"
	"ghostlang/lexer"
)

type parser struct {
	tokens []lexer.Token
	pos    int
}

func createParser(tokens []lexer.Token) *parser {
	return &parser{
		tokens,
		0,
	}
}
func Parse(tokens []lexer.Token) ast.BlockStmt {
	Body := make([]ast.Stmt, 0)
	p := createParser(tokens)

	for p.hasTokens() {
		Body = append(Body, parse_stmt(p))
	}

	return ast.BlockStmt{
		Body,
	}
}

func (p *parser) currentToken() lexer.Token {
	createTokenLookups()
	createTokenTypeLookups()
	return p.tokens[p.pos]
}

func (p *parser) advance() lexer.Token {
	token := p.currentToken()
	p.pos++
	return token
}

func (p *parser) nextToken() lexer.Token {
	return p.tokens[p.pos+1]
}

func (p *parser) hasTokens() bool {
	return p.pos < len(p.tokens) && p.currentTokenKind() != lexer.EOF
}

func (p *parser) currentTokenKind() lexer.TokenKind {
	return p.currentToken().Kind
}

func (p *parser) expectError(expectedKind lexer.TokenKind, err any) lexer.Token {
	token := p.currentToken()
	kind := token.Kind

	if kind != expectedKind {
		if err == nil {
			err = fmt.Sprintf("Expected token: %s but got %s\n", lexer.TypeString(expectedKind), lexer.TypeString(kind))
		}

		panic(err)
	}

	return p.advance()
}

func (p *parser) expect(expectedKind lexer.TokenKind) lexer.Token {
	return p.expectError(expectedKind, nil)
}
