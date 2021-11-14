package knowledge

import (
	"fmt"
	"log"
	"strings"
)

const (
	STATE_START           = 0
	STATE_EVENT           = 1
	STATE_EVENT_NAME_OPEN = 2
)

func getEffects(exprs []string) ([]*Relationship, []string, error) {
	ee := []*Relationship{}
	aa := []string{}
	for _, expr := range exprs {
		name, weight, ok, err := parseEventFunction(expr)
		if err != nil {
			return nil, nil, err
		}
		if ok {
			r := Relationship{
				Indirect: name,
				Weight:   weight,
				Effect:   nil,
			}
			ee = append(ee, &r)
		} else {
			aa = append(aa, expr)
		}
	}
	return ee, aa, nil
}

func parseEvent(name string, source string) (Event, string, error) {
	var err error
	e := NewEvent(name)
	snippet, source, err := getBlock(clean(source))
	if err != nil {
		return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
	}
	for _, token := range []string{"premises", "if", "then"} {
		if strings.HasPrefix(snippet, token) {
			snippet = snippet[len(token):]
			var block string
			var exprs []string
			block, snippet, err = getBlock(snippet)
			if err != nil {
				log.Print(snippet)
				return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
			}
			exprs, err = getExpressions(block)
			if err != nil {
				return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
			}
			switch token {
			case "premises":
				err := e.AddPremises(exprs)
				if err != nil {
					return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
				}
			case "if":
				e.AddConditions(exprs)
				if err != nil {
					return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
				}
			case "then":
				ee, aa, err := getEffects(exprs)
				if err != nil {
					return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
				}
				e.AddEffects(ee)
				err = e.AddAssignments(aa)
				if err != nil {
					return e, source, fmt.Errorf("failed parsing event %v -> %v", name, err)
				}
			}
		} else {
			return e, source, fmt.Errorf("failed parsing event %v -> expected '%v'", name, token)
		}
	}
	return e, source, nil
}

func Parse(source string) ([]*Event, error) {
	var err error
	ee := []*Event{}
	source = clean(source)
	for {
		l0 := len(source)
		if l0 == 0 {
			return ee, nil
		}
		if strings.HasPrefix(source, "/*") {
			_, source, err = readUntil(source, "*/")
			if err != nil {
				return ee, err
			}
		}
		if strings.HasPrefix(source, "event(") {
			var name string
			name, source, err = readUntil(source[6:], ")")
			if err != nil {
				return ee, err
			}
			var e Event
			e, source, err = parseEvent(name, source)
			if err != nil {
				return ee, err
			}
			ee = append(ee, &e)
		}
		if len(source) == l0 {
			return ee, fmt.Errorf("malformat source at char 0, size [%v]", l0)
		}
	}
}