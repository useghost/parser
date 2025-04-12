package parser

import (
	"ghostlang/ast"
	"ghostlang/lexer"
)

func parse_stmt(p *parser) ast.Stmt {
	stmt_fn, exists := stmt_lu[p.currentTokenKind()]

	if exists {
		return stmt_fn(p)
	}

	expression := parse_expr(p, default_bp)
	p.expect(lexer.SEMI_COLON)

	return ast.ExpressionStmt{Expression: expression}
}

func parse_declaration_stmt(p *parser) ast.Stmt {
	var explicitType ast.Type
	var assignedValue ast.Expr
	isConst := p.advance().Kind == lexer.CONST
	varName := p.expectError(lexer.IDENTIFIER, "Variable name identifier expected in variable declaration.")

	if p.currentTokenKind() == lexer.LESS {
		p.advance()
		explicitType = parse_type(p, default_bp)
		p.advance()
	}

	if p.currentTokenKind() != lexer.SEMI_COLON {
		p.expect(lexer.ASSINGMENT_EQUALS)
		assignedValue = parse_expr(p, assignment)
	} else if explicitType == nil {
		panic("Missing type or initial value in var declaration.")
	}

	if isConst && assignedValue == nil {
		panic("No assigned value in constant declaration.")
	}

	p.expect(lexer.SEMI_COLON)
	return ast.DeclarationStmt{
		IsConstant:    isConst,
		Identifier:    varName.Value,
		AssignedValue: assignedValue,
		ExplicitType:  explicitType,
	}
}
