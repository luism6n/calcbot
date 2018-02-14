package calc

import (
	"strconv"
	"testing"
)

func TestLexer(t *testing.T) {
	t.Run("Should recognize numbers", func(t *testing.T) {
		for _, number := range []string{"1", "2", "999", "100001"} {
			val, err := strconv.ParseFloat(number, 64)
			if err != nil {
				t.Fatalf("Something wrong with strconv.ParseFloat(%s, 64): %s", number, err.Error())
			}

			lexer := newCalcLexer(number)
			lval := &yySymType{}
			token := lexer.Lex(lval)

			if token != NUMBER || lval.val != val {
				t.Fatalf("%d != NUMBER or  %f != %f", token, lval.val, val)
			}
		}
	})
}
