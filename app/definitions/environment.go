package definitions

import "fmt"

type Environment struct {
	values map[string]any
}

func NewEnvironment() Environment {
	return Environment{
		values: make(map[string]any),
	}
}

func (env *Environment) Define(name string, value any) {
	env.values[name] = value
}

func (env *Environment) Get(name Token) (any, error) {
	val, ok := env.values[name.Lexeme]
	if ok {
		return val, nil
	}
	return nil, fmt.Errorf("Undefined variable '%s'.", name.Lexeme)
}

func (env *Environment) Assign(name Token, value any) error {
	if _, ok := env.values[name.Lexeme]; ok {
		env.values[name.Lexeme] = value
		return nil
	}
	return fmt.Errorf("undefined variable '%s'.", name.Lexeme)
}
