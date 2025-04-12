package parser

import (
	"fmt"
	"ghostlang/ast"
	"ghostlang/lexer"
)

type type_nud_handler func(p *parser) ast.Type
type type_led_handler func(p *parser, left ast.Type, bp bindingpower) ast.Type

type type_nud_lookup map[lexer.TokenKind]type_nud_handler
type type_led_lookup map[lexer.TokenKind]type_led_handler
type type_bp_lookup map[lexer.TokenKind]bindingpower

var type_bp_lu = bp_lookup{}
var type_nud_lu = type_nud_lookup{}
var type_led_lu = type_led_lookup{}

func type_led(kind lexer.TokenKind, bp bindingpower, led_fn type_led_handler) {
	type_bp_lu[kind] = bp
	type_led_lu[kind] = led_fn
}

func type_nud(kind lexer.TokenKind, nud_fn type_nud_handler) {
	type_nud_lu[kind] = nud_fn
}

func createTokenTypeLookups() {
	type_nud(lexer.IDENTIFIER, parse_symbol_type)
	type_nud(lexer.LEFT_BRACKET, parse_array_type)
}

func parse_symbol_type(p *parser) ast.Type {
	return ast.SymbolType{
		Name: p.expect(lexer.IDENTIFIER).Value,
	}
}

func parse_array_type(p *parser) ast.Type {
	p.advance()
	p.expect(lexer.RIGHT_BRACKET)
	underlyingType := parse_type(p, default_bp)
	return ast.ArrayType{
		UnderlyingType: underlyingType,
	}
}

func parse_type(p *parser, bp bindingpower) ast.Type {
	//parse nud
	tokenKind := p.currentTokenKind()
	nud_fn, exists := type_nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("type nud handler expected for token %s\n", lexer.TypeString(tokenKind)))
	}

	left := nud_fn(p)

	for type_bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := type_led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("type led handler expected for token %s\n", lexer.TypeString(tokenKind)))
		}
		left = led_fn(p, left, type_bp_lu[p.currentTokenKind()])
	}

	return left
}
