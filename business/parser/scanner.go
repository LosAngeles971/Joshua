package parser

import (
	"bytes"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
)

// List of all types of recognized tokens
const (
	text_token          = "text"    // arbitrary text
	if_token            = "if"      // keyword
	then_token          = "then"    // keyword
	effect_token        = "effects" // keyword
	open_event_token    = "event("  // keyword
	close_bracket_token = ")"       // keyword
	open_block_token    = "{"       // keyword
	close_block_token   = "}"       // keyword
	open_comment_token  = "/*"      // keyword
	close_comment_token = "*/"      // keyword
	quote_token         = "\""      // keyword
)

// List of all recognized keywords
var keywords = []string{
	if_token,
	then_token,
	effect_token,
	open_event_token,
	close_bracket_token,
	open_block_token,
	close_block_token,
	open_comment_token,
	close_comment_token,
	quote_token,
}

// List of chars which can be discarded
var trimmables = []rune{' ', '\t', '\n', '\r'}

// End of file char
var eof = rune(0)

// textStopTokens provides the possible stop tokens when the scanner is getting an arbitrary text.
// Indeed, the stop tokens can vary depending on the previous token of the arbitrary text.
var textStopTokens = map[string][]string{
	open_comment_token: {close_comment_token},
	open_event_token:   {close_bracket_token},
	quote_token:        {quote_token},
}

// trimAfterToken indicates if after a given token is possible to trim the trashable chars
var trimAfterToken = map[string]bool{
	text_token:          false,
	if_token:            true,
	then_token:          true,
	effect_token:        true,
	open_event_token:    false,
	close_bracket_token: true,
	open_block_token:    true,
	close_block_token:   true,
	open_comment_token:  false,
	close_comment_token: true,
	quote_token:         true,
}

// Token declares id and value of recognized tokens.
// Note that value is used only for text tokens.
type Token struct {
	id    string
	value string
}

// Scanner is a lexical scanner.
type Scanner struct {
	source []rune
	index  int
}

func NewScanner(source string) (*Scanner, error) {
	if len(source) < 1 {
		return nil, fmt.Errorf("source cannot be of size %v", len(source))
	}
	return &Scanner{
		source: []rune(source),
		index:  0,
	}, nil
}

// getAvailableSize provides the number chars not scanned yet.
func (s *Scanner) getAvailableSize() int {
	return len(s.source) - s.index
}

// eof says if the source code is completely scanned or not.
func (s *Scanner) eof() bool {
	return s.index >= len(s.source)
}

func (s *Scanner) getChar() rune {
	if s.index == len(s.source) {
		return eof
	}
	r := s.source[s.index]
	s.index++
	return r
}

func isTrimmable(c rune) bool {
	for i := range trimmables {
		if c == trimmables[i] {
			return true
		}
	}
	return false
}

// trim discards all continuous trimmable chars.
func (s *Scanner) trim() {
	for {
		if s.index >= len(s.source) {
			return
		}
		r := s.source[s.index]
		if isTrimmable(r) {
			s.index++
		} else {
			return
		}
	}
}

// getSubstring returns a substring from the source code.
func (s *Scanner) getSubstring(start int, lenght int) ([]rune, error) {
	if start >= len(s.source) || start+lenght-1 >= len(s.source) {
		return nil, fmt.Errorf("cannot get substring at [%v][%v] for a source of size %v", start, lenght, len(s.source))
	}
	extract := []rune{}
	for i := start; i < start+lenght; i++ {
		extract = append(extract, s.source[i])
	}
	return extract, nil
}

func isKeyword(expectedToken string) bool {
	for i := range keywords {
		if strings.EqualFold(keywords[i], expectedToken) {
			return true
		}
	}
	return false
}

// isNextKeyword says if next chars compose a recognized keyword or not.
func (s *Scanner) isNextKeyword(keyword string, move bool) (bool, error) {
	if !isKeyword(keyword) {
		return false, fmt.Errorf("expected keyword [%v] is not valid", keyword)
	}
	if len(keyword) > s.getAvailableSize() {
		return false, nil
	}
	extract, err := s.getSubstring(s.index, len(keyword))
	if err != nil {
		return false, err
	}
	log.Tracef("check if [%s] is equal to expected [%s]", string(extract), keyword)
	if strings.EqualFold(string(extract), keyword) {
		if move {
			s.index += len(keyword)
		}
		return true, nil
	}
	return false, nil
}

// getText extracts a token of arbitrary text until a recognized stop token.
func (s *Scanner) getText(stopTokens []string) (string, error) {
	var extract bytes.Buffer
	for {
		for i := range stopTokens {
			yes, err := s.isNextKeyword(stopTokens[i], false)
			if err != nil || yes {
				return extract.String(), err
			}
		}
		c := s.getChar()
		if c == eof {
			return extract.String(), nil
		}
		extract.WriteRune(c)
	}
}

// run scans the source code and returns a Lexer object.
func (s *Scanner) run() (*Lexer, error) {
	lexer := &Lexer{
		index:  0,
		tokens: []Token{},
	}
	var previous Token
	for {
		if s.eof() {
			return lexer, nil
		}
		if lexer.isEmpty() {
			s.trim()
		} else {
			previous, _ = lexer.last()
			if trimAfterToken[previous.id] {
				s.trim()
			}
		}
		changed := false
		// keyword
		for k := range keywords {
			yes, err := s.isNextKeyword(keywords[k], false)
			if err != nil {
				return nil, err
			}
			if yes {
				s.isNextKeyword(keywords[k], true)
				lexer.tokens = append(lexer.tokens, Token{
					id: keywords[k],
				})
				changed = true
			}
		}
		if !changed {
			// free text
			var stopTokens []string
			if lexer.isEmpty() {
				stopTokens = []string{quote_token}
			} else {
				var ok bool
				stopTokens, ok = textStopTokens[previous.id]
				if !ok {
					return nil, fmt.Errorf("free text cannot be after token id [%v]", previous.id)
				}
			}
			text, err := s.getText(stopTokens)
			if err != nil {
				return nil, err
			}
			lexer.tokens = append(lexer.tokens, Token{
				id:    text_token,
				value: text,
			})
		}
	}
}
