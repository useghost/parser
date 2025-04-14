package lexer

import (
	"fmt"
	"regexp"
)

type regexPattern struct {
	regex   *regexp.Regexp
	handler regexHandler
}

type lexer struct {
	patterns []regexPattern
	Tokens   []Token
	source   string
	pos      int
	line     int
}

func Tokenize(source string) []Token {
	lex := CreateLexer(source)

	for !lex.end_of_file() {
		matchfound := false
		for _, pattern := range lex.patterns {
			// fmt.Println(lex.patterns)
			loc := pattern.regex.FindStringIndex(lex.remainingSource())
			if loc != nil && loc[0] == 0 {
				pattern.handler(lex, pattern.regex)
				matchfound = true
				break
			}
		}

		if !matchfound {
			panic(fmt.Sprintf("unrecognized token at %v", lex.remainingSource))
		}
	}

	lex.push(NewToken(EOF, "EOF"))
	return lex.Tokens
}

func (lex *lexer) advanceBy(n int) {
	lex.pos += n
}

func (lex *lexer) currentByte() byte {
	return lex.source[lex.pos]
}

func (lex *lexer) advance() {
	lex.pos += 1
}

func (lex *lexer) remainingSource() string {
	return lex.source[lex.pos:]
}

func (lex *lexer) push(token Token) {
	lex.Tokens = append(lex.Tokens, token)
}

func (lex *lexer) end_of_file() bool {
	return lex.pos >= len(lex.source)
}

func CreateLexer(source string) *lexer {
	return &lexer{
		pos:    0,
		line:   1,
		source: source,
		Tokens: make([]Token, 0),
		patterns: []regexPattern{
			{regexp.MustCompile(`\s+`), skipHandler},
			{regexp.MustCompile(`\/\/.*`), commentHandler},
			{regexp.MustCompile(`"[^"]*"`), stringHandler},
			{regexp.MustCompile(`[+-]?[0-9]+(\.[0-9]+)?`), numberHandler},
			{regexp.MustCompile(`[a-zA-Z_][a-zA-Z0-9_]*`), symbolHandler},
			{regexp.MustCompile(`\[`), defaultHandler(LEFT_BRACKET, "[")},
			{regexp.MustCompile(`\]`), defaultHandler(RIGHT_BRACKET, "]")},
			{regexp.MustCompile(`\{`), defaultHandler(LEFT_BRACE, "{")},
			{regexp.MustCompile(`\}`), defaultHandler(RIGHT_BRACE, "}")},
			{regexp.MustCompile(`\(`), defaultHandler(LEFT_PAREN, "(")},
			{regexp.MustCompile(`\)`), defaultHandler(RIGHT_PAREN, ")")},
			{regexp.MustCompile(`==`), defaultHandler(EQUALS, "==")},
			{regexp.MustCompile(`!=`), defaultHandler(NOT_EQUALS, "!=")},
			{regexp.MustCompile(`=`), defaultHandler(ASSINGMENT_EQUALS, "=")},
			{regexp.MustCompile(`!`), defaultHandler(NOT, "!")},
			{regexp.MustCompile(`<=`), defaultHandler(LESS_EQUALS, "<=")},
			{regexp.MustCompile(`<`), defaultHandler(LESS, "<")},
			{regexp.MustCompile(`>=`), defaultHandler(GREATER_EQUALS, ">=")},
			{regexp.MustCompile(`>`), defaultHandler(GREATER, ">")},
			{regexp.MustCompile(`\|\|`), defaultHandler(OR, "||")},
			{regexp.MustCompile(`&&`), defaultHandler(AND, "&&")},
			{regexp.MustCompile(`\.\.`), defaultHandler(DOUBLE_DOT, "..")},
			{regexp.MustCompile(`\.`), defaultHandler(DOT, ".")},
			{regexp.MustCompile(`;`), defaultHandler(SEMI_COLON, ";")},
			{regexp.MustCompile(`:`), defaultHandler(COLON, ":")},
			{regexp.MustCompile(`\?`), defaultHandler(QUESTION_OPERATOR, "?")},
			{regexp.MustCompile(`,`), defaultHandler(COMMA, ",")},
			{regexp.MustCompile(`\+\+`), defaultHandler(PLUS_PLUS, "++")},
			{regexp.MustCompile(`--`), defaultHandler(MINUS_MINUS, "--")},
			{regexp.MustCompile(`\+=`), defaultHandler(PLUS_EQUALS, "+=")},
			{regexp.MustCompile(`-=`), defaultHandler(MINUS_EQUALS, "-=")},
			{regexp.MustCompile(`\+`), defaultHandler(PLUS, "+")},
			{regexp.MustCompile(`-`), defaultHandler(MINUS, "-")},
			{regexp.MustCompile(`/`), defaultHandler(DIVIDE, "/")},
			{regexp.MustCompile(`\*`), defaultHandler(MUL, "*")},
			{regexp.MustCompile(`%`), defaultHandler(MODULO, "%")},
		},
	}
}
func skipHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainingSource())
	lex.advanceBy(match[1])
}

func stringHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainingSource())
	stringLiteral := lex.remainingSource()[match[0]:match[1]]

	lex.push(NewToken(STRING, stringLiteral))
	lex.advanceBy(len(stringLiteral))
}

func numberHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainingSource())
	lex.push(NewToken(NUMBER, match))
	lex.advanceBy(len(match))
}

func symbolHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindString(lex.remainingSource())

	if kind, found := reserved_lu[match]; found {
		lex.push(NewToken(kind, match))
	} else {
		lex.push(NewToken(IDENTIFIER, match))
	}

	lex.advanceBy(len(match))
}

func defaultHandler(kind TokenKind, value string) regexHandler {
	return func(lex *lexer, _ *regexp.Regexp) {
		lex.advanceBy(len(value))
		lex.push(NewToken(kind, value))
	}
}

func commentHandler(lex *lexer, regex *regexp.Regexp) {
	match := regex.FindStringIndex(lex.remainingSource())
	if match != nil {
		// Advance past the entire comment.
		lex.advanceBy(match[1])
		lex.line++
	}
}

type regexHandler func(lex *lexer, regex *regexp.Regexp)
