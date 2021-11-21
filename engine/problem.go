package engine

import (
	"it/losangeles971/joshua/state"

	"gopkg.in/yaml.v3"
)

type Variable struct {
	Name string					`yaml:"name"`
	Value interface{}   		`yaml:"value"`
}

type Problem struct {
	Variables	[]Variable		`yaml:"variables"`
	Success		string			`yaml:"success"`
}

func LoadProblem(problemFile string) (state.State, string, error) {
	s := state.NewSimpleState()
	p := Problem{}
	err := yaml.Unmarshal([]byte(problemFile), &p)
	if err != nil {
		return nil, "", err
	}
	for _, v := range p.Variables {
		s.Add(v.Name, v.Value)
	}
	return s, p.Success, nil
}