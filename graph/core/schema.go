package core

import (
	"fmt"

	"github.com/google/uuid"
)

type URI struct {
	Namespace  string
	Identifier string
}

func NewURI(namespace string, id string) URI {
	return URI{
		Namespace: namespace,
		Identifier: id,
	}
}

func NewID(namespace string) URI {
	return URI{
		Namespace:  namespace,
		Identifier: uuid.NewString(),
	}
}

func (u URI) Clone() URI {
	return URI{
		Namespace:  u.Namespace,
		Identifier: u.Identifier,
	}
}

func (u URI) GetName() string {
	return fmt.Sprintf("%s::%s", u.Namespace, u.Identifier)
}

func (u1 URI) IsEqual(u2 URI) bool {
	if u1.Namespace == u2.Namespace && u1.Identifier == u2.Identifier {
		return true
	} else {
		return false
	}
}
