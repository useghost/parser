package parser

import (
	"ghostlang/ast"
	"ghostlang/lexer"
)

type bindingpower int

const (
	default_bp = iota
	comma
	assignment
	logical
	relational
	additive
	multiplicative
	unary
	call
	member
	primary
)

type stmt_handler func(p *parser) ast.Stmt
type nud_handler func(p *parser) ast.Expr
type led_handler func(p *parser, left ast.Expr, bp bindingpower) ast.Expr

type stmt_lookup map[lexer.TokenKind]stmt_handler
type nud_lookup map[lexer.TokenKind]nud_handler
type led_lookup map[lexer.TokenKind]led_handler
type bp_lookup map[lexer.TokenKind]bindingpower

var bp_lu = bp_lookup{}
var nud_lu = nud_lookup{}
var led_lu = led_lookup{}
var stmt_lu = stmt_lookup{}

func led(kind lexer.TokenKind, bp bindingpower, led_fn led_handler) {
	bp_lu[kind] = bp
	led_lu[kind] = led_fn
}

func nud(kind lexer.TokenKind, nud_fn nud_handler) {
	nud_lu[kind] = nud_fn
}

func stmt(kind lexer.TokenKind, stmt_fn stmt_handler) {
	bp_lu[kind] = default_bp
	stmt_lu[kind] = stmt_fn
}

func createTokenLookups() {
	led(lexer.ASSINGMENT_EQUALS, assignment, parse_assignment_expr)
	led(lexer.PLUS_EQUALS, assignment, parse_assignment_expr)
	led(lexer.MINUS_EQUALS, assignment, parse_assignment_expr)

	// Logical
	led(lexer.AND, logical, parse_binary_expr)
	led(lexer.OR, logical, parse_binary_expr)
	led(lexer.DOUBLE_DOT, logical, parse_binary_expr) //low precedence

	//member
	led(lexer.DOT, member, parse_member_expr)
	//call
	led(lexer.LEFT_PAREN, call, parse_call_expr)

	//Relational
	led(lexer.LESS, relational, parse_binary_expr)
	led(lexer.LESS_EQUALS, relational, parse_binary_expr)
	led(lexer.GREATER, relational, parse_binary_expr)
	led(lexer.GREATER_EQUALS, relational, parse_binary_expr)
	led(lexer.EQUALS, relational, parse_binary_expr)
	led(lexer.NOT_EQUALS, relational, parse_binary_expr)

	// < additive & multiplicative >
	led(lexer.PLUS, additive, parse_binary_expr)
	led(lexer.MINUS, additive, parse_binary_expr)
	led(lexer.MUL, additive, parse_binary_expr)
	led(lexer.DIVIDE, additive, parse_binary_expr)
	led(lexer.MODULO, additive, parse_binary_expr)

	//idfk
	nud(lexer.FN, parse_fn_expr)

	// < lit or symbol >
	nud(lexer.NUMBER, parse_primary_expr)
	nud(lexer.STRING, parse_primary_expr)
	nud(lexer.IDENTIFIER, parse_primary_expr)
	nud(lexer.MINUS, parse_prefix_expr)
	nud(lexer.LEFT_PAREN, parse_grouping_expr)
	nud(lexer.TYPENAME, parse_typename_expr)

	// < statements >

	stmt(lexer.CONST, parse_declaration_stmt)
	stmt(lexer.LET, parse_declaration_stmt)
	stmt(lexer.SET, parse_declaration_stmt)
	stmt(lexer.INFER, parse_infer_declaration_stmt)
	stmt(lexer.COMPILER, parse_compiler_option_statement)
	stmt(lexer.EXCLUDE, parse_exclude_fn_stmt)
	stmt(lexer.FN, parse_fn_declaration)
	stmt(lexer.RETURN, parse_return_stmt)
}
