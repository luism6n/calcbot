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
}

func floatEquals(a, b, eps float64) bool {
	if a >= b && a-b < eps {
		return true
	} else if a < b && b-a < eps {
		return true
	}

	return false
}
