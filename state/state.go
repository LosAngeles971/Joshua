/*
From Joshua's perspective, a state is a set of variables, the latter can be numbers of booleans.

 - This package defines the interface to handle a state
 - This package provides an implementation of the state's interface
*/
package state

const (
	TYPE_NUMBER = 0
	TYPE_BOOL   = 1
)

type State interface {
	Add(name string, value interface{}) error    // add a new variable into the state
	Update(name string, value interface{}) error // change the value of an existing variable
	Get(name string) (interface{}, bool)         // get the value of a variable
	Clone() State                                // return a clone of the state
	IsDefined(name string) bool                  // check is a variable exists
	GetType(name string) (int, error)            // return the type of a variable
	Vars() []string                              // return the list of the names of defined variables
	IsSubsetOf(father State) bool                // check if a state is a subset of another
	Translate() map[string]interface{}           // convert a state into a map of interface
}
