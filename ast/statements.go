package ast

type BlockStmt struct {
	Body []Stmt
}

func (n BlockStmt) stmt() {}

type ExpressionStmt struct {
	Expression Expr
}

func (n ExpressionStmt) stmt() {}

type DeclarationStmt struct {
	Identifier    string
	IsConstant    bool
	AssignedValue Expr
	ExplicitType  Type
}

func (n DeclarationStmt) stmt() {}
