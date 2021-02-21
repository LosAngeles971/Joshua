package problems

import (
	"io/ioutil"
	ctx "it/losangeles971/joshua/internal/context"
	"gopkg.in/yaml.v2"
)

type Variable struct {
	Name string					`yaml:"name"`
	Value float64				`yaml:"value"`
	Defined	bool				`yaml:"defined"`
	Range string 				`yaml:"range"`
}

type Problem struct {
	Variables	[]Variable		`yaml:"variables"`
	Success		string			`yaml:"success"`
}

func Load(problemFile string) (ctx.State, string, error) {
	c := ctx.CreateEmptyState()
	p := Problem{}
	b, err := ioutil.ReadFile(problemFile)
	if err != nil {
		return c, "", err
	}
	err = yaml.Unmarshal(b, &p)
	if err != nil {
		return c, "", err
	}
	for _, v := range p.Variables {
		vv := ctx.Variable{
			Name: v.Name,
			Defined: v.Defined,
			Value: v.Value,
		}
		if len(v.Range) > 0 {
			vv.Range, err = ctx.ParseRange(v.Range)
			if err != nil {
				return c, "", err
			}
		}
		c.Add(&vv)
	}
	return c, p.Success, nil
}