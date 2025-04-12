package ast

type SymbolType struct {
	Name string
}

func (t SymbolType) _type() {}

type ArrayType struct {
	UnderlyingType Type //[]int32
}

func (t ArrayType) _type() {}
