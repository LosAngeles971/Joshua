package knowledge

import (
	"errors"
	"fmt"
	ctx "it/losangeles971/joshua/internal/context"
	"math/rand"
)

type Person struct {
    DNA         map[string]float64
    locked      map[string]bool
    ranking     float64
}

func (p *Person) clone() *Person {
    pp := Person{
        DNA: p.DNA,
        locked: p.locked,
        ranking: p.ranking,
    }
    return &pp
}

func (p *Person) setState(s *ctx.State) error {
    for _, k := range s.Keys() {
        v, ok := p.DNA[k]
        if !ok {
            return errors.New("Missing variable into DNS: " + k)
        }
        s.Data[k].Defined = true
        s.Data[k].Value = v
    }
    return nil
}

type Population struct  {
    Fitness             *Path
    Size                int
    GeneticCode         ctx.State
    Population          []*Person
    Selective_pressure  float64     // usually between 1 and 2, greater value privileges good item for recombination
    Selection_ratio     float64     // usually 0.5, the size (%) of population to be replaced
    Mutation_rate       float64     // the probability of having one gene mutated inside a child
}

func (p Population) randomPerson() *Person {
    return p.Population[rand.Intn(len(p.Population))]
}

func (p *Population) Init() {
    p.Population = []*Person{}
    for i := 0; i < p.Size; i++ {
        person := Person{DNA: map[string]float64{}, locked: map[string]bool{}}
        for _, key := range p.GeneticCode.Keys() {
            v, _ := p.GeneticCode.Get(key)
            person.DNA[key] = v.GetRandom()
            person.locked[key] = false
        }
        p.Population = append(p.Population, &person)
    }
}

/*

*/
func (p *Population) ranking() error {
    for _, person := range p.Population {
        state := ctx.CreateEmptyState()
        for k, v := range person.DNA {
            state.Update(k, v)
        }
        p.Fitness.Reset()
        err := p.Fitness.Run(state, 0)
        if err != nil {
            return err
        }
        fmt.Println("Outcome: ", p.Fitness.GetOutcome())
        switch p.Fitness.GetOutcome() {
        case CE_OUTCOME_TRUE:
            for k := range person.DNA {
                person.locked[k] = true
            }
            person.ranking = 1.0
        case CE_OUTCOME_EFFECT_FALSE:
            person.ranking = 0.75
        case CE_OUTCOME_CAUSE_FALSE:
            person.ranking = 0.25
            for k := range person.DNA {
                person.locked[k] = false
            }
        default:
            person.ranking = 0
            for k := range person.DNA {
                person.locked[k] = false
            }
        }
    }
    return nil
}

func (p *Population) kill(i int) {
    // swap
    p.Population[len(p.Population) - 1], p.Population[i] = p.Population[i], p.Population[len(p.Population) -1]
    // reduce
    p.Population = p.Population[:len(p.Population) - 1]
}

/*
https://en.wikipedia.org/wiki/Fitness_proportionate_selection
*/
func (p Population) choose() int {
    total := 0.0
    for _, i := range p.Population {
        total += i.ranking
    }
    value := rand.Float64() * total
    for ix, i := range p.Population {
        value -= i.ranking
        if value <= 0 {
            return ix
        }
    }
    return len(p.Population) - 1
}

func (p *Population) selection() {
    alpha := []*Person{}
    selected := int(p.Selection_ratio * float64(len(p.Population)))
    for n := 0; n < selected; n++ {
        ix := p.choose()
        alpha = append(alpha, p.Population[ix].clone())
        p.kill(ix)
    }
    p.Population = alpha
 }

func (p Population) mutate(child *Person) {
    unlocked := []string{}
    for k := range child.DNA {
        if !child.locked[k] {
            unlocked = append(unlocked, k)
        }
    }
    if len(unlocked) == 0 {
        return
    }
    r := rand.Intn(len(unlocked))
	v, _ := p.GeneticCode.Get(unlocked[r])
	child.DNA[unlocked[r]] = v.GetRandom()
}

func (p Population) combine(mother *Person, father *Person) *Person {
    child := Person{}
    child.DNA = map[string]float64{} 
    child.locked = map[string]bool{}
    for key := range mother.DNA {
        if mother.locked[key] {
            child.DNA[key] = mother.DNA[key]
            child.locked[key] = true
        } else if father.locked[key] {
            child.DNA[key] = father.DNA[key]
            child.locked[key] = true
        } else {
            child.locked[key] = false
            r := rand.Float64()
            if r < 0.5 {
                child.DNA[key] = mother.DNA[key]
            } else {
                child.DNA[key] = father.DNA[key]
            }
        }
    }
    r := rand.Float64()
    if r < p.Mutation_rate {
        p.mutate(&child)
    }
    return &child
}

func (p *Population) crossover() {
    childs := []*Person{}
    for (len(p.Population) + len(childs)) < p.Size {
        mother := p.randomPerson()
        father := p.randomPerson()
        childs = append(childs, p.combine(mother, father))
    }
    p.Population = append(p.Population, childs...)
}

func (p *Population) GetOneSolution(s *ctx.State) (bool, error) {
    for _, person := range p.Population {
        if person.ranking == 1.0 {
            err := person.setState(s)
            if err != nil {
                return false, err
            } else {
                return true, nil
            }
        }
    }
    return false, nil
}

func MakePopulation(path *Path, init ctx.State) Population {
    return Population{
		Fitness: path,
		Size: 10,
		Selective_pressure: 1.8,
		Selection_ratio: 0.5,
		Mutation_rate: 0.75,
		GeneticCode: init,
	}
}

func CycleGenerations(p *Population, generations int) error {
    p.Init()
    for g := 0; g < generations; g++ {
        err := p.ranking()
        if err != nil {
            return err
        }
        p.selection()
        p.crossover()
    }
    return nil
}

