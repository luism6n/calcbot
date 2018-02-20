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
				t.Fatalf("%d != %d or  %f != %f", token, NUMBER, lval.val, c.Value)
			}
		}
	})

	t.Run("Should recognize identifiers", func(t *testing.T) {
		testCases := []string{
			"a",
			"aa",
			"a_b",
			"AB",
			"a1",
			"a",
		}

		for _, c := range testCases {
			lexer := newCalcLexer(c)
			lval := &yySymType{}
			token := lexer.Lex(lval)

			if token != IDENTIFIER || lval.name != c {
				t.Fatalf("%d != %d or  %q != %q", token, IDENTIFIER, lval.name, c)
			}
		}
	})
}

func floatEquals(a, b, eps float64) bool {
	if a >= b && a-b < eps {
		return true
	} else if a < b && b-a < eps {
		return true
	}

	return false
}
