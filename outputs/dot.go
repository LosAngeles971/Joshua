package outputs

import (
	"fmt"
	"io/ioutil"
	"it/losangeles971/joshua/knowledge"
)

func renderEdge(e *knowledge.Edge, cycle int) string {
	var label string 
	var color string
	switch e.Outcome {
	case knowledge.EFFECT_OUTCOME_CAUSE_FALSE:
		label = "c"
		color = "black"
	case knowledge.EFFECT_OUTCOME_EFFECT_FALSE:
		label = "e"
		color = "black"
	case knowledge.EFFECT_OUTCOME_ERROR:
		label = "?"
		color = "red"
	case knowledge.EFFECT_OUTCOME_LOOP:
		label = "l"
		color = "black"
	case knowledge.EFFECT_OUTCOME_NULL:
		label = "?"
		color = "red"
	case knowledge.EFFECT_OUTCOME_TRUE:
		label = "!"
		color = "green"
	case knowledge.EFFECT_OUTCOME_UNKNOWN:
		label = "?"
		color = "black"
	default:
		label = "#"
		color = "black"
	}
	return fmt.Sprintf("    %v -> %v [color=\"%s\",label=\"%d%s\"]",e.Cause.GetID(), e.Effect.GetID(), color, cycle, label)
}

func renderPath(p *knowledge.Path) (string, bool) {
	if !p.Executed {
		return "", false
	}
	output := ""
	ok := false
	switch p.Outcome {
	case knowledge.EFFECT_OUTCOME_CAUSE_FALSE:
	case knowledge.EFFECT_OUTCOME_EFFECT_FALSE:
		ok = true
	case knowledge.EFFECT_OUTCOME_ERROR:
	case knowledge.EFFECT_OUTCOME_LOOP:
	case knowledge.EFFECT_OUTCOME_NULL:
	case knowledge.EFFECT_OUTCOME_TRUE:
		ok = true
	case knowledge.EFFECT_OUTCOME_UNKNOWN:
	}
	if !ok {
		return "", false
	}
	for _, e := range p.Path {
		output += renderEdge(e, p.Cycle) + "\n"
	}
	return output, true
}

func GetDot(q knowledge.Queue) string {
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

func SaveDot(q knowledge.Queue, filename string) error {
	return ioutil.WriteFile(filename, []byte(GetDot(q)), 0644)
}