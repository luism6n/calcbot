//go:generate goyacc calc.y
package calc

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"unicode"
	"unicode/utf8"
)

var result float64 // yyParse stores the end result here

// Evaluate takes a program and returns the value of it's last statement.
func Evaluate(program string) (float64, error) {
	if yyParse(newCalcLexer(program)) != 0 {
		return 0.0, errors.New("Failed to parse program")
	}

	return result, nil
}

// Lexer is a math expressions (plus variables) tokenizer.
type calcLexer struct {
	program string
	ts, te  int // current token is program[ts:te]
}

// NewLexer returns a new lexer for the given program.
func newCalcLexer(program string) *calcLexer {
	return &calcLexer{
		program: program,
		ts:      -1, // current token's start
		te:      0,  // and end positions
	}
}

// Lex returns the next token type and puts its value (if any) in lval.
func (l *calcLexer) Lex(lval *yySymType) int {
	l.consumeWhiteSpace()

	if l.eof() {
		return 0
	}

	reLog := regexp.MustCompile(`log`)
	rePow := regexp.MustCompile(`pow`)
	reOp := regexp.MustCompile(`[;=,()+/*-]`)
	reIdent := regexp.MustCompile(`\pL(\pL|[0-9_])*`)

	// This scary-looking regex was taken from
	// https://golang.org/ref/spec#Floating-point_literals
	// with the added option to have no decimal point. Funny enough, putting the
	// [0-9]+ at the beginning fails to match 1.5, for example.
	// TODO: find documentation about in what order Go tries to match the ORed
	// regexes.
	reNumber := regexp.MustCompile(`[0-9]+\.([0-9]+)?([eE][+-]?[0-9]+)?|[0-9]+([eE][+-]?[0-9]+)|\.[0-9]+([eE][+-]?[0-9]+)?|[0-9]+`)

	switch {
	case l.matchAndAdvance(reLog):
		return LOG
	case l.matchAndAdvance(rePow):
		return POW
	case l.matchAndAdvance(reOp):
		return int(l.currentToken()[0])
	case l.matchAndAdvance(reIdent):
		lval.name = l.currentToken()
		return IDENTIFIER
	case l.matchAndAdvance(reNumber):
		lval.val = l.parseFloat()
		return NUMBER
	default:
		l.Error(fmt.Sprintf("Error parsing expression: %s", l.program[l.te:]))
		return -1
	}
}

func (l *calcLexer) eof() bool {
	return l.te == len(l.program)
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

// matchAndAdvance will check if the beginning of the remainder of the program
// l.program matches the regular expression re. If it does, l.te is advanced by
// the size of the match and true is returned. Else, nothing happens and false
// is returned.
func (l *calcLexer) matchAndAdvance(re *regexp.Regexp) bool {
	if loc := re.FindStringIndex(l.program[l.te:]); loc != nil && loc[0] == 0 {
		l.ts, l.te = l.te, l.te+loc[1]
		return true
	}

	return false
}

func (l *calcLexer) parseFloat() float64 {
	val, err := strconv.ParseFloat(l.currentToken(), 64)
	if err != nil {
		l.Error(fmt.Sprintf("ParseFloat(%s, 64) failed: %s", l.currentToken(), err.Error()))
	}
	return val
}

func (l *calcLexer) currentToken() string {
	return l.program[l.ts:l.te]
}

// Error is called when something is wrong in the Lexer's program.
func (l *calcLexer) Error(s string) {
	log.Printf("Syntax error: %s\n", s)
}
