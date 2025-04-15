package parser

import (
	"fmt"
	"ghostlang/ast"
	"ghostlang/lexer"
	"strconv"
)

func parse_expr(p *parser, bp bindingpower) ast.Expr {
	//parse nud
	tokenKind := p.currentTokenKind()
	nud_fn, exists := nud_lu[tokenKind]

	if !exists {
		panic(fmt.Sprintf("nud handler expected for token %s\n", lexer.TypeString(tokenKind)))
	}

	left := nud_fn(p)

	for bp_lu[p.currentTokenKind()] > bp {
		tokenKind := p.currentTokenKind()
		led_fn, exists := led_lu[tokenKind]
		if !exists {
			panic(fmt.Sprintf("led handler expected for token %s\n", lexer.TypeString(tokenKind)))
		}
		left = led_fn(p, left, bp)
	}

	return left
}
func parse_primary_expr(p *parser) ast.Expr {
	switch p.currentTokenKind() {
	case lexer.NUMBER:
		value := p.advance().Value
		number, _ := strconv.ParseFloat(value, 64)
		// var numType string
		return ast.NumberExpr{
			Value:        value,
			Float64Value: number,
			// Type: inferType()
			// ParsedValue: parseNumber()
		}

	case lexer.STRING:
		return ast.StringExpr{
			Value: p.advance().Value,
		}
	case lexer.IDENTIFIER:
		return ast.SymbolExpr{
			Value: p.advance().Value,
		}
	default:
		panic(fmt.Sprintf("cant create primary expression from %s\n", lexer.TypeString(p.currentTokenKind())))
	}
}

func parse_member_expr(p *parser, left ast.Expr, bp bindingpower) ast.Expr {
	// Assume we've just seen a DOT and the current token is the property name
	p.advance() // consume the dot

	if p.currentTokenKind() != lexer.IDENTIFIER {
		panic(fmt.Sprintf("expected identifier after '.', got %s", lexer.TypeString(p.currentTokenKind())))
	}

	property := ast.SymbolExpr{
		Value: p.advance().Value,
	}

	return ast.MemberExpr{
		Object:   left,
		Property: property,
	}
}

func parse_binary_expr(p *parser, left ast.Expr, bp bindingpower) ast.Expr {
	operatorToken := p.advance()
	right := parse_expr(p, bp_lu[p.currentTokenKind()])

	return ast.BinaryExpr{
		Left:     left,
		Operator: operatorToken,
		Right:    right,
	}
}

func parse_prefix_expr(p *parser) ast.Expr {
	operatorToken := p.advance()
	rhs := parse_expr(p, default_bp)

	return ast.PrefixExpr{
		Operator:  operatorToken,
		RightExpr: rhs,
	}
}

func parse_typename_expr(p *parser) ast.Expr {
	p.advance()
	exp := parse_expr(p, default_bp)
	return ast.TypenameExpression{
		Expression: exp,
	}
}
func parse_assignment_expr(p *parser, left ast.Expr, bp bindingpower) ast.Expr {
	operatorToken := p.advance()
	rhs := parse_expr(p, assignment)
	return ast.AssignmentExpr{
		Operator: operatorToken,
		Value:    rhs,
		Assignee: left,
	}
}

func parse_grouping_expr(p *parser) ast.Expr {
	left := p.advance()
	expr := parse_expr(p, default_bp)
	right := p.expect(lexer.RIGHT_PAREN)
	return ast.GroupExpression{
		Opener:     left,
		Expression: expr,
		Closer:     right,
	}
}

func parse_call_expr(p *parser, left ast.Expr, bp bindingpower) ast.Expr {
	p.advance()
	arguments := make([]ast.Expr, 0)

	for p.hasTokens() && p.currentTokenKind() != lexer.RIGHT_PAREN {
		arguments = append(arguments, parse_expr(p, assignment))

		if p.currentTokenKind() != lexer.RIGHT_PAREN {
			p.expect(lexer.COMMA)
		}
	}

	p.expect(lexer.RIGHT_PAREN)
	return ast.CallExpr{
		Method:    left,
		Arguments: arguments,
	}
}

func parse_fn_expr_params_body(p *parser) ([]ast.Argument, ast.Type, []ast.Stmt) {
	p.expect(lexer.LEFT_PAREN)
	functionParams := make([]ast.Argument, 0)
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

func parse_fn_expr(p *parser) ast.Expr {
	p.advance()
	if !p.currentToken().IsOfTypes(lexer.IDENTIFIER, lexer.LEFT_PAREN) {
		panic("Error: Parser did not recieve either function name or left parenthesis in function expression")
	}
	var nameToken lexer.Token = p.currentToken()

	if nameToken.Kind == lexer.IDENTIFIER {
		p.advance()
		functionParams, returnType, functionBody := parse_fn_expr_params_body(p)
		return ast.AnonymousFunctionExpr{
			Arguments:    functionParams,
			ReturnType:   returnType,
			Body:         functionBody,
			OptionalName: nameToken.Value,
		}
	} else {
		functionParams, returnType, functionBody := parse_fn_expr_params_body(p)
		return ast.AnonymousFunctionExpr{
			Arguments:    functionParams,
			ReturnType:   returnType,
			Body:         functionBody,
			OptionalName: "",
		}
	}

}
