package calc

import (
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("Should recognize numbers", func(t *testing.T) {
		testCases := []struct {
			Input string
			Value float64
		}{
			{"1", 1.0},
			{"2", 2.0},
			{"999", 999.0},
			{"100001", 100001.0},
			{"1.5", 1.5},
			{".1", 0.1},
			{"0.", 0.},
			{"72.40", 72.40},
			{"072.40", 072.40},
			{"2.71828", 2.71828},
			{"1.e+0", 1.e+0},
			{"6.67428e-11", 6.67428e-11},
			{"1E6", 1E6},
			{".25", .25},
			{".12345E+5", .12345E+5},
			{"\t6.67428e-11", 6.67428e-11},
			{"6.67428e-11\t", 6.67428e-11},
			{"6.67428e-11 ", 6.67428e-11},
			{" 6.67428e-11", 6.67428e-11},
			{"\n6.67428e-11", 6.67428e-11},
			{"6.67428e-11\n", 6.67428e-11},
		}

		for _, c := range testCases {
			lexer := newCalcLexer(c.Input)
			lval := &yySymType{}
			token := lexer.Lex(lval)

			if token != NUMBER || !floatEquals(lval.val, c.Value, 0.000001) {
				t.Fatalf("%s != %s or  %f != %f", tokname(token), tokname(NUMBER), lval.val, c.Value)
			}
		}
	})

	t.Run("Should recognize identifiers", func(t *testing.T) {
		testCases := []struct {
			Input string
			Value string
		}{
			{"a", "a"},
			{"aa", "aa"},
			{"a_b", "a_b"},
			{"AB", "AB"},
			{"a1", "a1"},
			{"a", "a"},
			{" a", "a"},
			{"a ", "a"},
			{"\ta", "a"},
			{"a\t", "a"},
			{"\na", "a"},
			{"a\n", "a"},
		}

		for _, c := range testCases {
			lexer := newCalcLexer(c.Input)
			lval := &yySymType{}
			token := lexer.Lex(lval)

			if token != IDENTIFIER || lval.name != c.Value {
				t.Fatalf("%s != %s or  %q != %q", tokname(token), tokname(IDENTIFIER), lval.name, c)
			}
		}
	})

	t.Run("Should recognize operations and parenthesis", func(t *testing.T) {
		testCases := []struct {
			Input     string
			TokenType int
		}{
			{"+", '+'},
			{"-", '-'},
			{"*", '*'},
			{"/", '/'},
			{"log", LOG},
			{"pow", POW},
		}

		for _, c := range testCases {
			lexer := newCalcLexer(c.Input)
			lval := &yySymType{}
			token := lexer.Lex(lval)

			if token != c.TokenType {
				t.Fatalf("%s != %s for %s", tokname(token), tokname(c.TokenType), c.Input)
			}
		}
	})
}

func tokname(token int) string {
	// This is probably a bad idea
	index := token - NUMBER + 3
	if index >= 0 && index < len(yyToknames) {
		return yyToknames[index]
	}

	return string([]byte{byte(token)})
}

func floatEquals(a, b, eps float64) bool {
	if a >= b && a-b < eps {
		return true
	} else if a < b && b-a < eps {
		return true
	}

	return false
}
