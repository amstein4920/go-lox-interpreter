package tree

type Class struct {
	Name string
}

func NewClass(name string) Class {
	return Class{
		Name: name,
	}
}

func (c Class) String() string {
	return c.Name
}

func (c Class) Call(interpreter Interpreter, arguments []any) any {
	return NewInstance(c)
}

func (c Class) Arity() int {
	return 0
}
