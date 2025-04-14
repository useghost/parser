package parser

import (
	"fmt"
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

func parse_infer_declaration_stmt(p *parser) ast.Stmt {
	p.advance()
	var assignedValue ast.Expr
	var varName string
	isConst := false

	if p.currentToken().Kind == lexer.CONST {
		isConst = true
		p.advance()
	}

	varName = p.currentToken().Value
	p.advance()

	if p.currentTokenKind() != lexer.SEMI_COLON {
		if p.currentTokenKind() == lexer.LESS {
			panic("Unexpected type declaration for inferred variable")
		}
		p.expect(lexer.ASSINGMENT_EQUALS)
		assignedValue = parse_expr(p, assignment)
	}

	if isConst && assignedValue == nil {
		panic("error: no assigned value in constant declaration.")
	}

	p.expect(lexer.SEMI_COLON)
	return ast.DeclarationStmt{
		IsConstant:    isConst,
		Identifier:    varName,
		AssignedValue: assignedValue,
		ExplicitType:  nil,
	}
}

func parse_compiler_option_statement(p *parser) ast.Stmt {
	p.expect(lexer.COMPILER)

	if p.currentTokenKind() != lexer.OPTION {
		panic(fmt.Sprintf("Error: expected keyword \"option\" after token compiler but got %s", lexer.TypeString(p.currentTokenKind())))
	}
	p.advance()

	varName := p.currentToken().Value
	p.advance()

	var assignedValue ast.Expr

	if p.currentTokenKind() == lexer.ASSINGMENT_EQUALS {
		p.advance()

		assignedValue = parse_expr(p, default_bp)
	}

	p.expect(lexer.SEMI_COLON)

	return ast.CompilerOptionStmt{
		OptionName:  varName,
		OptionValue: assignedValue,
	}
}

func parse_block_stmt(p *parser) ast.Stmt {
	p.expect(lexer.LEFT_BRACE)
	body := []ast.Stmt{}

	for p.hasTokens() && p.currentTokenKind() != lexer.RIGHT_BRACE {
		body = append(body, parse_stmt(p))
	}

	p.expect(lexer.RIGHT_BRACE)
	return ast.BlockStmt{
		Body: body,
	}
}

func parse_fn_params_and_body(p *parser) ([]ast.Argument, ast.Type, []ast.Stmt) {
	functionParams := make([]ast.Argument, 0)

	p.expect(lexer.LEFT_PAREN)
	for p.hasTokens() && p.currentTokenKind() != lexer.RIGHT_PAREN {
		paramName := p.expect(lexer.IDENTIFIER).Value
		paramType := parse_type(p, default_bp)

		functionParams = append(functionParams, ast.Argument{
			Name: paramName,
			Type: paramType,
		})

		if !p.currentToken().IsOfTypes(lexer.RIGHT_PAREN, lexer.EOF) {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.RIGHT_PAREN)
	var returnType ast.Type

	if p.currentTokenKind() != lexer.LEFT_BRACE {
		returnType = parse_type(p, default_bp)
	}

	functionBody := ast.ExpectStmt[ast.BlockStmt](parse_block_stmt(p)).Body

	return functionParams, returnType, functionBody
}

func parse_fn_declaration(p *parser) ast.Stmt {
	p.advance()
	functionName := p.expect(lexer.IDENTIFIER).Value
	functionParams, returnType, functionBody := parse_fn_params_and_body(p)

	return ast.FunctionDeclarationStmt{
		Arguments:  functionParams,
		ReturnType: returnType,
		Body:       functionBody,
		Name:       functionName,
	}
}
