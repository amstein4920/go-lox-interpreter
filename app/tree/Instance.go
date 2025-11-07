package tree

import (
	"errors"
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/app/definitions"
)

type Instance struct {
	klass  Class
	fields map[string]any
}

func NewInstance(klass Class) Instance {
	return Instance{
		klass:  klass,
		fields: make(map[string]any),
	}
}

func (i *Instance) get(name definitions.Token) (any, error) {
	if value, ok := i.fields[name.Lexeme]; ok {
		return value, nil
	}
	if method, ok := i.klass.findMethod(name.Lexeme); ok {
		return method, nil
	}
	return nil, errors.New("Undefined property")
}

func (i *Instance) set(name definitions.Token, value any) {
	i.fields[name.Lexeme] = value
}

func (i Instance) String() string {
	return fmt.Sprintf("%s instance", i.klass.Name)
}
