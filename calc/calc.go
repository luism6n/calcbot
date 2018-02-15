//go:generate goyacc calc.y
package calc

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// Lexer is a math expressions (plus variables) tokenizer.
type calcLexer struct {
	program string
	ts, te  int // current token is program[ts:te]
}

// NewLexer returns a new lexer for the given program.
func newCalcLexer(program string) *calcLexer {
	return &calcLexer{
		program: program,
		ts:      -1,
		te:      0,
	}
}

// Lex returns the next token type and puts its value (if any) in lval.
func (l *calcLexer) Lex(lval *yySymType) int {
	l.nextToken()

	var err error
	lval.val, err = strconv.ParseFloat(l.currentToken(), 64)
	if err != nil {
		l.Error(fmt.Sprintf("ParseFloat(%s, 64) failed: %s", l.currentToken(), err.Error()))
	}
	return NUMBER
}

func (l *calcLexer) nextToken() {
	re := regexp.MustCompile(`[0-9]+(\.[0-9]*)?`)
	loc := re.FindStringIndex(l.program[l.te:])
	l.ts, l.te = loc[0], loc[1]
}

func (l *calcLexer) eof() bool {
	return l.ts == len(l.program)
}

func (l *calcLexer) nextRune() rune {
	c, width := utf8.DecodeRuneInString(l.program[l.te:])
	l.te += width

	return c
}

func (l *calcLexer) currentToken() string {
	return l.program[l.ts:l.te]
}

// Error is called when something is wrong in the Lexer's program.
func (l *calcLexer) Error(s string) {
	log.Fatalf("Syntax error: %s", s)
}
