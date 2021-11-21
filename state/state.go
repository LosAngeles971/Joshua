package state

const (
	TYPE_NUMBER = 0
	TYPE_BOOL   = 1
)

type State interface {
	Add(name string, value interface{}) error
	Update(name string, value interface{}) error
	Get(name string) (interface{}, bool)
	Clone() State
	IsDefined(name string) bool
	GetType(name string) (int, error)
	Vars() []string
	IsSubsetOf(father State) bool
	Translate() map[string]interface{}
}
