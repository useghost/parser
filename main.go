package main

import (
	"ghostlang/lexer"
	"ghostlang/parser"
	"os"

	"github.com/sanity-io/litter"
)

func main() {
	sourceBytes, _ := os.ReadFile("vars.lang")
	source := string(sourceBytes)

	tokens := lexer.Tokenize(string(source))

	ast := parser.Parse(tokens)
	litter.Dump(ast)

}
