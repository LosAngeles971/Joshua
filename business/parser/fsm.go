package parser

import "fmt"

const (
	call_action  = "call"
	token_action = "token"
)

// Transition describes a possible transition from one state to the next one
type Transition struct {
	Action string `yaml:"action"` // the Action can be a "token" or a "call" to a sub-process based on another FSM
	Token  string `yaml:"token"`  // the expected token, if the Action is "token"
	Sub    string `yaml:"sub"`    // the sub-process name, if the Action is "call"
	Min    int    `yaml:"min"`    // the minimum occorrences of the sub-process, if the Action is "call"
}

// FSM is a really simple deterministic finite-state-machine
// there is one FSM for each Joshua's atomic block of code
type FSM struct {
	ID          string       `yaml:"id"`
	Transitions []Transition `yaml:"transitions"`
}

// firstToken returns the exepected token from the first transition
func (f FSM) firstToken() (string, error) {
	if len(f.Transitions) < 1 {
		return "", fmt.Errorf("FSM [%s] -> no transitions", f.ID)
	}
	if f.Transitions[0].Action != token_action {
		return "", fmt.Errorf("FSM [%s] -> first transition must be 'token' not '%s'", f.ID, f.Transitions[0].Action)
	}
	return f.Transitions[0].Token, nil
}

// isApplicable checks if the FSM is applicable for the incoming block
// at the date, this function requires the FSM's first transition must be a token not a call
func (f FSM) isApplicable(t Token) (bool, error) {
	tt, err := f.firstToken()
	if err != nil {
		return false, err
	}
	return tt == t.id, nil
}
