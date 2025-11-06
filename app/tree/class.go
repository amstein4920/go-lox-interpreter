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

func (c *Class) call(interpreter Interpreter, arguments []any) any {
	return NewInstance(*c)
}

func (c *Class) arity() int {
	return 0
}
