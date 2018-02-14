//go:generate goyacc calc.y
package calc

import "log"

// Lexer is a math expressions (plus variables) tokenizer.
type calcLexer struct {
	program string
	pos     int
}

// NewLexer returns a new lexer for the given program.
func newCalcLexer(program string) calcLexer {
	return calcLexer{
		program: program,
	}
}

// Lex returns the next token type and puts its value (if any) in lval.
func (l calcLexer) Lex(lval *yySymType) int {
	return NUMBER
}

// Error is called when something is wrong in the Lexer's program.
func (l calcLexer) Error(s string) {
	log.Fatalf("Syntax error: %s", s)
}
