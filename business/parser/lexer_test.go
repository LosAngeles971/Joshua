package parser

import (
	"strings"
	"testing"
)

type lexerTest struct {
	source   string
	tokens   []int
	literals []string
}

var lexerTests = []lexerTest{
	// {
	// 	source:   "  aa  if",
	// 	tokens:   []int{space_token, identifier_token, eof_token},
	// 	literals: []string{"  ", "aa  if", ""},
	// },
	{
		source:   " /* aa */ if ",
		tokens:   []int{space_token, open_comment_block, space_token, close_comment_block, space_token, if_token, space_token},
		literals: []string{" ", "/*", " ", "aa ", "*/", " ", "if", " "},
	},
}

func TestLexer(t *testing.T) {
	for i, test := range lexerTests {
		parser := NewScanner(strings.NewReader(test.source))
		tokens := []int{}
		literals := []string{}
		for {
			tt, ll := parser.Scan()
			tokens = append(tokens, tt)
			literals = append(literals, ll)
			if tt == eof_token {
				break
			}
		}
		if len(tokens) != len(test.tokens) || len(literals) != len(test.literals) || len(tokens) != len(literals) {
			t.Fatalf("test %v failed", i)
		}
		for j := range tokens {
			if tokens[j] != test.tokens[j] || literals[j] != test.literals[j] {
				t.Fatalf("test %v failed, expected [%v][%s] at %v not [%v][%s]", i, test.tokens[j], test.literals[j], j, tokens[j],literals[j])
			}
		}
	}
}