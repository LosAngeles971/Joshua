package genetic

import (
    "math/rand"
    "it/losangeles971/joshua/internal/runtime"
    "it/losangeles971/joshua/internal/knowledge"
    ctx "it/losangeles971/joshua/internal/context"
)

type Openrange struct {
    Min     float64
    Max     float64
}

func (g Openrange) Get() float64 {
    return g.Min + rand.Float64() * (g.Max - g.Min)
}

type Set struct {
    Values []float64
}

func (g Set) Get() float64 {
    return g.Values[rand.Intn(len(g.Values))]
}

type Gene interface {
    Get() float64
}

type Person struct {
    DNA         map[string]float64
    Fitness     float64
}

type Population struct  {
    Fitness             runtime.Path
    Size                int
    GeneticCode         map[string]Gene
    Population          []Person
    Selective_pressure  float64     // usually between 1 and 2, greater value privileges good item for recombination
    Selection_ratio     float64     // usually 0.5, the size (%) of population to be replaced
    Mutation_rate       float64     // the probability of having one gene mutated inside a child
}

func (p *Population) AddGene(key string, g Gene) {
    p.GeneticCode[key] = g
}

func (p Population) RandomPerson() Person {
    return p.Population[rand.Intn(len(p.Population))]
}

func (p *Population) Init() {
    p.Population = []Person{}
    for i := 0; i < p.Size; i++ {
        person := Person{DNA: map[string]float64{}}
        for key, gene := range p.GeneticCode {
            person.DNA[key] = gene.Get()
        }
        p.Population = append(p.Population, person)
    }
}

/*

*/
func (p *Population) ranking() error {
    for _, person := range p.Population {
        state := ctx.Create()
        for k, v := range person.DNA {
            state.Add(k, v)
        }
        err := p.Fitness.Run(state, 0)
        if err != nil {
            return err
        }
        switch p.Fitness.Outcome() {
        case knowledge.CE_OUTCOME_TRUE:
            person.Fitness = 1.0
        case knowledge.CE_OUTCOME_EFFECT_FALSE:
            person.Fitness = 0.75
        case knowledge.CE_OUTCOME_CAUSE_FALSE:
            person.Fitness = 0.25
        default:
            person.Fitness = 0
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
        total += i.Fitness
    }
    value := rand.Float64() * total
    for ix, i := range p.Population {
        value -= i.Fitness
        if value <= 0 {
            return ix
        }
    }
    return len(p.Population) - 1
}

func (p *Population) selection() {
    alpha := []Person{}
    for n := 0; n < int(p.Selection_ratio * float64(len(p.Population))); n++ {
        ix := p.choose()
        alpha = append(alpha, p.Population[ix])
        p.kill(ix)
    }
    p.Population = alpha
 }

func (p Population) mutate(child *Person) {
    // an approach to randomly select a key of the DNA
    r := rand.Intn(len(child.DNA))
	for k := range child.DNA {
		if r == 0 {
			child.DNA[k] = p.GeneticCode[k].Get()
            return
		}
		r--
	}
	panic("unreachable")
}

func (p Population) combine(mother Person, father Person) Person {
    child := Person{}
    child.DNA = map[string]float64{} 
    for key := range mother.DNA {
        r := rand.Float64()
        if r < 0.5 {
            child.DNA[key] = mother.DNA[key]
        } else {
            child.DNA[key] = father.DNA[key]
        }
    }
    r := rand.Float64()
    if r < p.Mutation_rate {
        p.mutate(&child)
    }
    return child
}

func (p *Population) crossover() {
    childs := []Person{}
    for (len(p.Population) + len(childs)) < p.Size {
        mother := p.RandomPerson()
        father := p.RandomPerson()
        childs = append(childs, p.combine(mother, father))
    }
    p.Population = append(p.Population, childs...)
}

func Cycle(p *Population, generations int) {
    p.Init()
    for g := 0; g < generations; g++ {
        p.ranking()
        p.selection()
        p.crossover()
    }
}