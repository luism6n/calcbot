package calc

import "testing"

func TestLexer(t *testing.T) {
	t.Run("Should recognize numbers", func(t *testing.T) {
		lexer := newCalcLexer("1")
		token := lexer.Lex(nil)
		if token != NUMBER {
			t.Fatalf("%d != NUMBER", token)
		}
	})
}
