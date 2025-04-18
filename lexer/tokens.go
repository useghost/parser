package lexer

import "fmt"

type TokenKind int

// some must be non case-sensitive, match it in regex.
const (
	EOF TokenKind = iota
	NULL
	TRUE
	FALSE
	NUMBER
	STRING
	IDENTIFIER

	//GROUPING
	LEFT_BRACKET
	RIGHT_BRACKET
	LEFT_BRACE
	RIGHT_BRACE
	LEFT_PAREN
	RIGHT_PAREN

	//Equality
	ASSINGMENT_EQUALS
	EQUALS // equality check; ==
	NOT_EQUALS
	NOT

	//Continued equivalence
	LESS
	LESS_EQUALS
	GREATER
	GREATER_EQUALS

	//Logical
	OR
	AND

	//Smybols
	DOT
	DOUBLE_DOT      //.. - [0..10]
	SPREAD_OPERATOR //...
	SEMI_COLON
	COLON
	QUESTION_OPERATOR
	COMMA

	//shorthand operators
	PLUS_PLUS
	MINUS_MINUS
	PLUS_EQUALS
	MINUS_EQUALS

	//Math operators
	PLUS
	MINUS
	DIVIDE
	MUL
	MODULO

	//Reserved keywords
	LET //maybe change to set?
	SET
	INFER
	CONST

	COMPILER
	OPTION

	IMPORT
	AS

	CLASS
	NEW
	FROM //import {a,b,c,} from "test.g" or import wholemodule from "test.g"
	FN   //funcdef
	RETURN
	IF
	ELSE
	ELSEIF
	FOR
	FOR_EACH
	WHILE
	EXPORT
	AT_OPERATOR //@something above functions,
	WASM_EXPORT //WASM.Export in the form of @WASM.Export
	INCLUDES
	TYPENAME //typename "hi" ==
	STRUCT

	//MISC
	EXCLUDE //exclude fn (is excluded when compiling, not just not exported)
)

var reserved_lu map[string]TokenKind = map[string]TokenKind{
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"let":      LET,
	"set":      SET,
	"infer":    INFER,
	"const":    CONST,
	"class":    CLASS,
	"option":   OPTION,
	"exclude":  EXCLUDE,
	"new":      NEW,
	"import":   IMPORT,
	"from":     FROM,
	"fn":       FN,
	"as":       AS,
	"return":   RETURN,
	"if":       IF,
	"else":     ELSE,
	"compiler": COMPILER,
	"foreach":  FOR_EACH,
	"while":    WHILE,
	"for":      FOR,
	"export":   EXPORT,
	"typename": TYPENAME,
	"includes": INCLUDES,
	"struct":   STRUCT,
}

var TypeStrings = []string{
	"EOF",
	"NULL",
	"TRUE",
	"FALSE",
	"NUMBER",
	"STRING",
	"IDENTIFIER",
	"LEFT_BRACKET",
	"RIGHT_BRACKET",
	"LEFT_BRACE",
	"RIGHT_BRACE",
	"LEFT_PAREN",
	"RIGHT_PAREN",
	"ASSINGMENT_EQUALS",
	"EQUALS",
	"NOT_EQUALS",
	"NOT",
	"LESS",
	"LESS_EQUALS",
	"GREATER",
	"GREATER_EQUALS",
	"OR",
	"AND",

	"DOT",
	"DOUBLE_DOT",
	"SPREAD_OPERATOR",
	"SEMI_COLON",
	"COLON",
	"QUESTION_OPERATOR",
	"COMMA",

	"PLUS_PLUS",
	"MINUS_MINUS",
	"PLUS_EQUALS",
	"MINUS_EQUALS",
	"PLUS",
	"MINUS",
	"DIVIDE",
	"MUL",
	"MODULO",

	"LET",
	"SET",
	"INFER",
	"CONST",

	"COMPILER",
	"OPTION",
	"IMPORT",
	"AS",
	"CLASS",
	"NEW",
	"FROM",
	"FN",
	"RETURN",
	"IF",
	"ELSE",
	"ELSEIF",
	"FOR",
	"FOR_EACH",
	"WHILE",
	"EXPORT",
	"AT_OPERATOR",
	"WASM_EXPORT",
	"INCLUDES",
	"TYPENAME",
	"STRUCT",
	"EXCLUDE",
}

type Token struct {
	Kind  TokenKind
	Value string
}

func (token Token) IsOfTypes(tokenTypes ...TokenKind) bool {
	for _, tokenType := range tokenTypes {
		if tokenType == token.Kind {
			return true
		}
	}

	return false
}

func (token Token) Print() {
	if token.IsOfTypes(IDENTIFIER, NUMBER, STRING) {
		fmt.Printf("%s(%s)\n", TypeString(token.Kind), token.Value)
	} else {
		fmt.Printf("%s()\n", TypeString(token.Kind))
	}
}
func TypeString(kind TokenKind) string {
	return TypeStrings[int(kind)]
}

func NewToken(kind TokenKind, value string) Token {
	return Token{
		kind, value,
	}
}
