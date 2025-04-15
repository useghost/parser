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

type MemberExpr struct {
	Object   Expr       // could be SymbolExpr, or nested MemberExpr
	Property SymbolExpr // the right-hand identifier (e.g. 'name' in 'user.name')
}

func (m MemberExpr) expr() {}

type CallExpr struct {
	Method    Expr
	Arguments []Expr
}

func (n CallExpr) expr() {}

type AnonymousFunctionExpr struct {
	Arguments    []Argument
	Body         []Stmt
	OptionalName string
	ReturnType   Type
}

func (n AnonymousFunctionExpr) expr() {}

type ByteStringExpr struct {
	Value string
}

func (n ByteStringExpr) expr() {}

type TypenameExpression struct {
	Expression Expr
}

func (n TypenameExpression) expr() {}
