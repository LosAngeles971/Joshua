package parser

import "fmt"

// Lexer is the result of a source's scanning
// It is a ordinated list of tokens and it provides functions for browsing
type Lexer struct {
	tokens []Token
	index  int
}

func (l *Lexer) last() (Token, error) {
	if len(l.tokens) == 0 {
		return Token{}, fmt.Errorf("empy Lexer")
	}
	return l.tokens[len(l.tokens)-1], nil
}

func (l *Lexer) isEmpty() bool {
	return len(l.tokens) == 0
}

func (l *Lexer) isDrained() bool {
	return l.index >= len(l.tokens)
}

func (l *Lexer) readToken() (Token, error) {
	if l.index >= len(l.tokens) {
		return Token{}, fmt.Errorf("reach the end of lexer")
	}
	return l.tokens[l.index], nil
}

func (l *Lexer) getToken() (Token, error) {
	t, err := l.readToken()
	if err != nil {
		return Token{}, err
	}
	l.index++
	return t, nil
}