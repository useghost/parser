package ast

import "ghostlang/helpers"

type Stmt interface {
	stmt()
}

type Expr interface {
	expr()
}

type Type interface {
	_type()
}

func ExpectStmt[T Stmt](expr Stmt) T {
	return helpers.ExpectType[T](expr)
}
