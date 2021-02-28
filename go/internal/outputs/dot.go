package outputs

import (
	kkk "it/losangeles971/joshua/internal/knowledge"
	"io/ioutil"
	"fmt"
)

func renderEdge(e *kkk.Edge, cycle int) string {
	var label string 
	var color string
	switch e.Outcome {
	case kkk.CE_OUTCOME_CAUSE_FALSE:
		label = "c"
		color = "black"
	case kkk.CE_OUTCOME_EFFECT_FALSE:
		label = "e"
		color = "black"
	case kkk.CE_OUTCOME_ERROR:
		label = "?"
		color = "red"
	case kkk.CE_OUTCOME_LOOP:
		label = "l"
		color = "black"
	case kkk.CE_OUTCOME_NULL:
		label = "?"
		color = "red"
	case kkk.CE_OUTCOME_TRUE:
		label = "!"
		color = "green"
	case kkk.CE_OUTCOME_UNKNOWN:
		label = "?"
		color = "black"
	default:
		label = "#"
		color = "black"
	}
	return fmt.Sprintf("    %d -> %d [color=\"%s\",label=\"%d%s\"]",e.Cause.UID, e.Effect.UID, color, cycle, label)
}

func renderPath(p *kkk.Path) (string, bool) {
	if !p.Executed {
		return "", false
	}
	output := ""
	ok := false
	switch p.Outcome {
	case kkk.CE_OUTCOME_CAUSE_FALSE:
	case kkk.CE_OUTCOME_EFFECT_FALSE:
		ok = true
	case kkk.CE_OUTCOME_ERROR:
	case kkk.CE_OUTCOME_LOOP:
	case kkk.CE_OUTCOME_NULL:
	case kkk.CE_OUTCOME_TRUE:
		ok = true
	case kkk.CE_OUTCOME_UNKNOWN:
	}
	if !ok {
		return "", false
	}
	for _, e := range p.Path {
		output += renderEdge(e, p.Cycle) + "\n"
	}
	return output, true
}

func GetDot(q kkk.Queue) string {
	output := "digraph solution {\n    node [shape = circle];\n"
	for _, p := range q.Paths {
		if p.Cycle != -1 {
			render, ok := renderPath(p)
			if ok {
				output += render
			}
		}
	}
	output += "}"
	return output
}

func SaveDot(q kkk.Queue, filename string) error {
	return ioutil.WriteFile(filename, []byte(GetDot(q)), 0644)
}