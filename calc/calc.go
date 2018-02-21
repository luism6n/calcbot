//go:generate goyacc calc.y
package calc

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"unicode"
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
	if l.eof() {
		return 0
	}

	l.consumeWhiteSpace()

	reLog := regexp.MustCompile(`log`)
	rePow := regexp.MustCompile(`pow`)
	reOp := regexp.MustCompile(`[+/*-]`)
	reIdent := regexp.MustCompile(`[a-zA-Z][_a-zA-Z0-9]*`)

	// This scary-looking regex was taken from
	// https://golang.org/ref/spec#Floating-point_literals
	// with the added option to have no decimal point. Funny enough, putting the
	// [0-9]+ at the beginning fails to match 1.5, for example.
	// TODO: find documentation about in what order Go tries to match the ORed
	// regexes.
	reNumber := regexp.MustCompile(`[0-9]+\.([0-9]+)?([eE][+-]?[0-9]+)?|[0-9]+([eE][+-]?[0-9]+)|\.[0-9]+([eE][+-]?[0-9]+)?|[0-9]+`)

	var tokenType int
	var loc []int
	if loc = reLog.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te+loc[0], l.te+loc[1]
		tokenType = LOG
	} else if loc = rePow.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te+loc[0], l.te+loc[1]
		tokenType = POW
	} else if loc = reOp.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te+loc[0], l.te+loc[1]
		tokenType = int(l.currentToken()[0])
	} else if loc = reIdent.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te+loc[0], l.te+loc[1]
		tokenType = IDENTIFIER
		lval.name = l.currentToken()
	} else if loc = reNumber.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te+loc[0], l.te+loc[1]
		tokenType = NUMBER
		var err error
		lval.val, err = strconv.ParseFloat(l.currentToken(), 64)
		if err != nil {
			l.Error(fmt.Sprintf("ParseFloat(%s, 64) failed: %s", l.currentToken(), err.Error()))
		}
	}

	return tokenType
}

func (l *calcLexer) eof() bool {
	return l.ts == len(l.program)
}

func (l *calcLexer) consumeWhiteSpace() {
	c := l.peekRune()
	for unicode.IsSpace(c) {
		l.nextRune()
		c = l.peekRune()
	}
}

func (l *calcLexer) peekRune() rune {
	c, _ := utf8.DecodeRuneInString(l.program[l.te:])
	return c
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
