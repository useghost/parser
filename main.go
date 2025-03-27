package main

import (
	"os"
	"ghostlang/lexer"
)

func main() {
	sourceBytes, _ := os.ReadFile("test.lang")
	source := string(sourceBytes)
	
	tks := lexer.Tokenize(string(source))
	for _, token := range tks {
		token.Print()
	}
}