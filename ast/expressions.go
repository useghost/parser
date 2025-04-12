package ast

import "ghostlang/lexer"

// literal expressions
type NumberExpr struct {
	Value        string
	Float64Value float64
}

func (n NumberExpr) expr() {}

type StringExpr struct {
	Value string
}

func (n StringExpr) expr() {}

type SymbolExpr struct {
	Value string
}

func (n SymbolExpr) expr() {}

// complex expressions

type BinaryExpr struct {
	Left     Expr
	Operator lexer.Token
	Right    Expr
}

func (n BinaryExpr) expr() {}

type PrefixExpr struct {
	Operator  lexer.Token
	RightExpr Expr
}

func (n PrefixExpr) expr() {}

type AssignmentExpr struct {
	Assignee Expr
	Operator lexer.Token
	Value    Expr
}

func (n AssignmentExpr) expr() {}

type GroupExpression struct {
	Opener     lexer.Token
	Expression Expr
	Closer     lexer.Token
}

func (n GroupExpression) expr() {}
