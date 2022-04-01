package parser

import (
	"fmt"
	"it/losangeles971/joshua/business/knowledge"
	"it/losangeles971/joshua/business/math"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
)

type Parser struct {
	s *Scanner
}

func NewParser(source string) *Parser {
	return &Parser{s: NewScanner(strings.NewReader(source))}
}

func (p *Parser) ignoreSpace() (int, string) {
	token, literal := p.s.Scan()
	if token == space_token {
		token, literal = p.s.Scan()
	}
	return token, literal
}

func (p *Parser) ignoreComments() (int, string) {
	token, literal := p.s.Scan()
	if token != close_comment_block {
		token, literal = p.s.Scan()
	}
	return token, literal
}

func (p *Parser) ignoreSpaceAndComments() (int, string) {
	for {
		token, text := p.ignoreSpace()
		switch token {
		case eof_token:
			return eof_token, ""
		case open_comment_block:
			p.ignoreComments()
		default:
			return token, text
		}
	}
}

func (p *Parser) parseExpressions() ([]*govaluate.EvaluableExpression, error) {
	exprs := []*govaluate.EvaluableExpression{}
	for {
		token, expr := p.ignoreSpaceAndComments()
		if token == eof_token || token == close_block_token {
			return exprs, nil
		}
		if token == identifier_token {
			e, err := math.NewExpression(expr)
			if err != nil {
				return nil, fmt.Errorf("invalid expression [%s] -> [%v]", expr, err)
			}
			exprs = append(exprs, e)
		}
		token, _ = p.ignoreSpaceAndComments() // must be a ;
		if token != close_expression {
			return nil, fmt.Errorf("expected [%s]", keywords[close_expression])
		}
	}
}

func (p *Parser) parseRelationship() (*knowledge.Relationship, error) {
	token, name := p.ignoreSpace()
	if token != identifier_token {
		return nil, fmt.Errorf("expected relationship identifier")
	}
	token, _ = p.ignoreSpace()
	if token != comma_token {
		return nil, fmt.Errorf("expected [%s]", keywords[comma_token])
	} 
	token, weight := p.ignoreSpace()
	if token != identifier_token {
		return nil, fmt.Errorf("expected number")
	}
	w, err := strconv.ParseFloat(weight, 64)
	if err != nil {
		return nil, fmt.Errorf("expected number")
	}
	token, _ = p.ignoreSpace()
	if token != close_event_token {
		return nil, fmt.Errorf("expected [%s]", keywords[comma_token])
	}
	return knowledge.NewRelationship(name, knowledge.WithWeight(w)), nil
}

func (p *Parser) parseThen() ([]knowledge.Assignment, []*knowledge.Relationship, error) {
	exprs := []knowledge.Assignment{}
	rel := []*knowledge.Relationship{}
	for {
		token, literal := p.ignoreSpaceAndComments()
		if token == eof_token || token == close_block_token {
			return exprs, rel, nil
		}
		if token == open_event_token {
			r, err := p.parseRelationship()
			if err != nil {
				return nil, nil, err
			}
			rel = append(rel, r)
		} else if token == identifier_token {
			a, err := knowledge.NewAssignment(literal)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid expression [%s] -> [%v]", literal, err)
			}
			exprs = append(exprs, a)
		} else {
			return nil, nil, fmt.Errorf("expected [%s] or assignment expression", keywords[open_event_token])
		}
		token, _ = p.ignoreSpaceAndComments() // must be a ;
		if token != close_expression {
			return nil, nil, fmt.Errorf("expected [%s]", keywords[close_expression])
		}
	}
}

func (p *Parser) parseEvent() (*knowledge.Event, error) {
	token, _ := p.ignoreSpaceAndComments() // must be "event("
	if token == eof_token {
		return nil, nil
	}
	if token != open_event_token {
		return nil, fmt.Errorf("expected open event [%s]", keywords[open_event_token])
	}
	token, name := p.s.scanToken()
	if token != identifier_token {
		return nil, fmt.Errorf("invalid name for event -> [%s]", name)
	}
	token, _ = p.ignoreSpaceAndComments() // must be ")"
	if token != close_event_token {
		return nil, fmt.Errorf("expected closing event name [%s]", keywords[close_event_token])
	}
	token, _ = p.ignoreSpaceAndComments() // must be "{"
	if token != open_block_token {
		return nil, fmt.Errorf("expected open block [%s]", keywords[open_block_token])
	}
	token, _ = p.ignoreSpaceAndComments() // must be "if"
	if token != if_token {
		return nil, fmt.Errorf("expected [%s]", keywords[if_token])
	}
	token, _ = p.ignoreSpaceAndComments() // must be "{"
	if token != open_block_token {
		return nil, fmt.Errorf("expected [%s]", keywords[open_block_token])
	}
	if_exprs, err := p.parseExpressions()
	if err != nil {
		return nil, err
	}
	token, _ = p.ignoreSpaceAndComments() // must be "then"
	if token != then_token {
		return nil, fmt.Errorf("expected [%s]", keywords[then_token])
	}
	token, _ = p.ignoreSpaceAndComments() // must be "{"
	if token != open_block_token {
		return nil, fmt.Errorf("expected [%s]", keywords[open_block_token])
	}
	let_exprs, rel, err := p.parseThen()
	if err != nil {
		return nil, err
	}
	return knowledge.NewEvent(
		name, 
		knowledge.WithAssignments(let_exprs), 
		knowledge.WithConditions(if_exprs),
		knowledge.WithRelationships(rel)), nil
}

func (p *Parser) Parse() ([]*knowledge.Event, error) {
	events := []*knowledge.Event{}
	for {
		e, err := p.parseEvent()
		if err != nil {
			return nil, err
		}
		if e == nil {
			return events, nil
		}
		events = append(events, e)
	}
}