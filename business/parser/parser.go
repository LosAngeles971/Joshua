package parser

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed joshua.yml
var joshua string

type Parser struct {
	Lib   []FSM `yaml:"automata"`
	lexer *Lexer
	code  []EventCode
}

func NewParser(l *Lexer) (*Parser, error) {
	if l == nil {
		return nil, fmt.Errorf("lexer cannot be nil")
	}
	fsm := &Parser{
		lexer: l,
		code:  []EventCode{},
	}
	err := yaml.Unmarshal([]byte(joshua), fsm)
	return fsm, err
}

func (p *Parser) setCodeName(name string) {
	p.code = append(p.code, EventCode{name: name})
}

func (p *Parser) addIf(expr string) {
	e := p.code[len(p.code)-1]
	e.conditions = append(e.conditions, expr)
}

func (p *Parser) addThen(expr string) {
	e := p.code[len(p.code)-1]
	e.assignments = append(e.assignments, expr)
}

func (p *Parser) addEffect(expr string) {
	e := p.code[len(p.code)-1]
	e.effects = append(e.effects, expr)
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
					p.addIf(token.value)
				}
			case "thenexpression":
				if token.id == text_token {
					p.addThen(token.value)
				}
			case "thenevent":
				if token.id == text_token {
					p.addEffect(token.value)
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
