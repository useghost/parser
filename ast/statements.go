package ast

type Argument struct {
	Name string
	Type Type
}

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

type CompilerOptionStmt struct {
	OptionName  string
	OptionValue Expr
}

func (n CompilerOptionStmt) stmt() {}

type FunctionDeclarationStmt struct {
	Arguments  []Argument
	Name       string
	Body       []Stmt
	ReturnType Type
}

func (n FunctionDeclarationStmt) stmt() {}

type ReturnStatement struct {
	ValueExpression Expr
}

func (n ReturnStatement) stmt() {}