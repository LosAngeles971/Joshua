package parser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const (
	illegal_token       = 0
	space_token         = 1
	if_token            = 3
	then_token          = 4
	identifier_token    = 5
	eof_token           = 6
	open_event_token    = 7
	close_event_token   = 8
	open_block_token    = 9
	close_block_token   = 10
	open_comment_block  = 11
	close_comment_block = 12
	close_expression    = 13
	comma_token         = 14

	state_init = "init"
	state_event = "event"
)

var keywords = map[int]string{
	if_token:            "IF",
	then_token:          "THEN",
	open_event_token:    "EVENT(",
	close_event_token:   ")",
	open_block_token:    "{",
	close_block_token:   "}",
	open_comment_block:  "/*",
	close_comment_block: "*/",
	close_expression:    ";",
	comma_token:         ",",
}

type fsmState struct {
	id       string
	transitions map[int]string
	stop     []int
}

var fsm = []fsmState{
	{
		id: state_init,
		transitions: map[int]string{
			open_event_token: state_event,
		},
		stop: []int{eof_token},
	},
}

var eof = rune(0)

// Lexical scanner.
type Scanner struct {
	state string
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		state: state_init,
		r: bufio.NewReader(r),
	}
}

// read reads the next rune from the bufferred reader.
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// undo returns a previously read rune back.
func (s *Scanner) undo() {
	_ = s.r.UnreadRune()
}

func (s *Scanner) isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func (s *Scanner) isStopChar(ch rune) bool {
	switch ch {
	case eof, '*', ')':
		return true
	}
	return false
}

func isKeyword(token string) (int, bool) {
	for k, literal := range keywords {
		if literal == strings.ToUpper(token) {
			return k, true
		}
	}
	return identifier_token, false
}

// scanSpace extracts all contiguous whitespace.
func (s *Scanner) scanSpace() (int, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(s.read())
	for {
		if ch := s.read(); ch == eof {
			return space_token, buffer.String()
		} else if !s.isSpace(ch) {
			s.undo()
			return space_token, buffer.String()
		} else {
			buffer.WriteRune(ch)
		}
	}
}

func (s *Scanner) scanToken() (int, string) {
	var buffer bytes.Buffer
	// first char must exist and it is not a space
	for {
		ch := s.read()
		if s.isStopChar(ch) {
			s.undo()
			token, ok := isKeyword(buffer.String())
			if ok {
				return token, buffer.String()
			}
			return identifier_token, buffer.String()
		} else {
			_, _ = buffer.WriteRune(ch)
			token, ok := isKeyword(buffer.String())
			if ok {
				return token, buffer.String()
			}
		}
	}
}

// Scan returns next token
func (s *Scanner) Scan() (int, string) {
	ch := s.read()
	switch ch {
	case eof:
		return eof_token, ""
	case ' ':
		s.undo()
		return s.scanSpace()
	default:
		s.undo()
		return s.scanToken()
	}
}
