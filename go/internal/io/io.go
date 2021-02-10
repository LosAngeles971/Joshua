package io

import (
	"errors"
	"io/ioutil"
	"crypto/rand"
	"os"
	"fmt"
	"gopkg.in/yaml.v2"
	"it/losangeles971/joshua/internal/problems"
	kkk "it/losangeles971/joshua/internal/knowledge"
	ctx "it/losangeles971/joshua/internal/context"
)

func getUUID() (string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		os.Exit(1)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// A link expresses a cause-effect relationship
type Link struct {
	Causes 	[]string	`yaml:"cause"`
	Effects	[]string	`yaml:"effect"`
	Type 	float64		`yaml:"type"`
}

type Storage struct {
	Events	[]kkk.Event	`yaml:"events"`
	Links  	[]Link	`yaml:"links"`
}

func Load(file interface{}) (kkk.Knowledge, error) {
	knowledge := kkk.Knowledge{}
	k := Storage{}
	var b []byte
	var err error
	switch file.(type) {
	case string:
		b, err = ioutil.ReadFile(file.(string))
		if err != nil {
			return knowledge, err
		}
	case []byte:
		b = file.([]byte)
	default:
		return knowledge, errors.New("Unrecognized type of input")
	}
	if err := yaml.Unmarshal(b, &k); err != nil {
		return knowledge, err
	}
	knowledge.Events = k.Events
	knowledge.Relationships = []kkk.Relationship{}
	for _, l := range k.Links {
		var cause kkk.Event
		var ok bool
		if len(l.Causes) > 1{
			cause = kkk.Event{
				ID: getUUID(),
				Dependecies: []kkk.Event{},
			}
			for _, f := range l.Causes {
				c, ok := knowledge.GetEvent(f)
				if !ok {
					return knowledge, errors.New("Event does not exist: " + f)
				}
				cause.Dependecies = append(cause.Dependecies, c)
			}
		} else if len(l.Causes) == 1 {
			cause, ok = knowledge.GetEvent(l.Causes[0])
			if !ok {
				return knowledge, errors.New("Event does not exist: " + l.Causes[0])
			}
		}
		for _, e := range l.Effects {
			effect, ok := knowledge.GetEvent(e)
			if !ok {
				return knowledge, errors.New("Event does not exist: " + e)
			}
			knowledge.Relationships = append(knowledge.Relationships, kkk.Relationship{
				Cause: cause,
				Effect: effect,
				Type: l.Type,
			})
		}
	}
	return knowledge, nil
}

// for testing purpose
func (k *Storage) Parse(yml interface{}) error {
	switch yml.(type) {
	case string:
		return yaml.Unmarshal([]byte(yml.(string)), k)
	case []byte:
		return yaml.Unmarshal(yml.([]byte), k)
	default:
		return errors.New("Unrecognized type of input")
	}
}

// for testing purpose
func (l *Link) Load(file string) error {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, l)
}

func LoadProblem(knowledgeFile string, problemFile string) (ctx.State, kkk.Knowledge, kkk.Event, error) {
	k, err := Load(knowledgeFile)
	if err != nil {
		return ctx.Create(), kkk.Knowledge{}, kkk.Event{}, err
	}
	init, success_name, err := problems.Load(problemFile)
	if err != nil {
		return ctx.Create(), k, kkk.Event{}, err
	}
	success, ok := k.GetEvent(success_name)
	if !ok {
		return init, k, kkk.Event{}, err
	}
	return init, k, success, nil
}