package parser

import (
	_ "embed"
	"testing"
)

//go:embed test1.joshua
var test1 string

func TestGetAvailableSize(t *testing.T) {
	source := "Neque porro quisquam est qui dolorem ipsum quia dolor sit amet"
	tests := map[int]int{
		0: 62,
		5: 57,
		62: 0,
	}
	s, err := NewScanner(source)
	if err != nil {
		t.Fatal(err)
	}
	for location, size := range tests {
		s.index = location
		if s.getAvailableSize() != size {
			t.Errorf("failed test with location %v and size %v not expected %v", location, size, s.getAvailableSize())
		}
	}
}

func TestIsNextToken(t *testing.T) {
	source := ")---if---then---}---event("
	tests := map[int]string{
		0: close_bracket_token,
		4: if_token,
		9: then_token,
		16: close_block_token,
		20: open_event_token,
	}
	s, err := NewScanner(source)
	if err != nil {
		t.Fatal(err)
	}
	for location, token := range tests {
		s.index = location
		yes, err := s.isNextKeyword(token, false)
		if err != nil {
			t.Fatal(err)
		}
		if !yes {
			t.Errorf("failed test at location %v with expected %v", location, token)
		}
	}
}

func TestGetText(t *testing.T) {
	source := ")---if---then---}---event("
	tests := map[int]string{
		0: "",
		1: "---",
	}
	s, err := NewScanner(source)
	if err != nil {
		t.Fatal(err)
	}
	for location, token := range tests {
		s.index = location
		tt, err := s.getText([]string{if_token, close_bracket_token})
		if err != nil {
			t.Fatal(err)
		}
		if tt != token {
			t.Errorf("failed test at location %v with expected %v not %v", location, token, tt)
		}
	}
}

func TestScanner(t *testing.T) {
	test_tokens := []Token{
		{ id: open_comment_token,},
		{ id: text_token, },
		{ id: close_comment_token,},
		{ id: open_event_token,	},
		{ id: text_token, },
		{ id: close_bracket_token, },
		{ id: open_block_token,},
		{ id: if_token,	},
		{ id: open_block_token,	},
		{ id: quote_token, }, { id: text_token, }, { id: quote_token, },
		{ id: quote_token, }, { id: text_token, }, { id: quote_token, },
		{ id: close_block_token, },
		{ id: then_token, },
		{ id: open_block_token,	},
		{ id: close_block_token, },
		{ id: close_block_token,},
	}
	s, err := NewScanner(test1)
	if err != nil {
		t.Fatal(err)
	}
	lexer, err := s.run()
	if err != nil {
		t.Fatal(err)
	}
	tt := lexer.tokens
	if len(tt) != len(test_tokens) {
			t.Errorf("expected %v tokens not %v", len(test_tokens), len(tt))
	}
	for j := range test_tokens {
		if j < len(tt) {
			if tt[j].id == text_token {
				t.Logf("[[%s]]", tt[j].value)
			}
			if tt[j].id != test_tokens[j].id {
				t.Errorf("[%v] ERROR expected [%v] token not [%v]", j, tt[j].id, test_tokens[j].id)
			} else {
				t.Logf("[%v] expected [%v] token found [%v]", j, tt[j].id, test_tokens[j].id)
			}
		}
	}
}