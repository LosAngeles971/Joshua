package parser

import (
	_ "embed"
	"fmt"
	"it/losangeles971/joshua/business/knowledge"
	"it/losangeles971/joshua/business/math"

	"gopkg.in/yaml.v3"
)

//go:embed joshua.yml
var joshua string

type Parser struct {
	Lib   []FSM `yaml:"automata"`
	lexer *Lexer
	code  []*knowledge.Event
}

func NewParser(l *Lexer) (*Parser, error) {
	if l == nil {
		return nil, fmt.Errorf("lexer cannot be nil")
	}
	fsm := &Parser{
		lexer: l,
		code: []*knowledge.Event{},
	}
	err := yaml.Unmarshal([]byte(joshua), fsm)
	return fsm, err
}

func (p *Parser) setCodeName(name string) {
	p.code = append(p.code, knowledge.NewEvent(name))
}

func (p *Parser) addIf(expr string) error {
	e := p.code[len(p.code)-1]
	cc, err := math.NewExpression(expr)
	if err != nil {
		return err
	}
	e.AddCondition(cc)
	return nil
}

func (p *Parser) addThen(expr string) error {
	e := p.code[len(p.code)-1]
	aa, err := knowledge.NewAssignment(expr)
	if err != nil {
		return err
	}
	e.AddAssigment(aa)
	return nil
}

func (p *Parser) addEffect(expr string) error {
	e := p.code[len(p.code)-1]
	e.AddEffect(knowledge.NewRelationship(expr))
	return nil
}

func (p *Parser) getFSM(id string) (FSM, error) {
	for i := range p.Lib {
		if p.Lib[i].ID == id {
			return p.Lib[i], nil
		}
	}
	return FSM{}, fmt.Errorf("FSM [%s] not defined", id)
}

func (p *Parser) process(fsmID string) error {
	fsm, err := p.getFSM(fsmID)
	if err != nil {
		return err
	}
	for i, trx := range fsm.Transitions {
		if trx.Action == call_action {
			sub, err := p.getFSM(trx.Sub)
			if err != nil {
				return err
			}
			count := 0
			keepgoing := true
			for keepgoing {
				if p.lexer.isDrained() {
					if count >= trx.Min {
						keepgoing = false
					} else {
						return fmt.Errorf("FSM [%s][%v] -> expected at least one [%v]", fsm.ID, i, trx.Sub)
					}
				} else {
					nextToken, err := p.lexer.readToken()
					if err != nil {
						return err
					}
					yes, err := sub.isApplicable(nextToken)
					if err != nil {
						return err
					}
					if yes {
						err := p.process(trx.Sub)
						if err != nil {
							return err
						}
						count++
					} else {
						if count >= trx.Min {
							keepgoing = false
						} else {
							return fmt.Errorf("FSM [%s][%v] -> expected at least one [%v]", fsm.ID, i, trx.Sub)
						}
					}
				}
			}
		} else {
			if p.lexer.isDrained() {
				return fmt.Errorf("FSM [%s][%v] -> expected [%v]", fsm.ID, i, trx.Token)
			}
			token, _ := p.lexer.getToken()
			if token.id != trx.Token {
				if i == 0 {
					return nil
				} else {
					return fmt.Errorf("FSM [%s][%v] -> expected [%v] not [%v]", fsm.ID, i, trx.Token, token.id)
				}
			}
			switch fsm.ID {
			case "event":
				if token.id == text_token {
					p.setCodeName(token.value)
				}
			case "ifexpression":
				if token.id == text_token {
					err := p.addIf(token.value)
					if err != nil {
						return err
					}
				}
			case "thenexpression":
				if token.id == text_token {
					err := p.addThen(token.value)
					if err != nil {
						return err
					}
				}
			case "thenevent":
				if token.id == text_token {
					err := p.addEffect(token.value)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (p *Parser) Parse() error {
	if p.lexer.isEmpty() {
		return nil
	}
	return p.process("body")
}
